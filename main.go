package main

import (
	"ChiHtmx/database"
	"ChiHtmx/middlewares"
	"ChiHtmx/model"
	"log"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	database.ConnectDB()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/user-info", userInfoHandler)

	r.Get("/posts", postHandler)

	r.Get("/post/create", createPostHandler)
	r.Post("/post/create", createPostHandler)

	r.Route("/post/{id}", func(r chi.Router) {
		r.Use(middlewares.PostCtx)

		// post object fetched in the PostCtx middleware. Handlers can perform its own specific set of actions.
		r.Get("/", getPostHandler)

		r.Get("/edit", editPostHandler)
		r.Post("/edit", editPostHandler)

	})

	http.ListenAndServe(":3000", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := make(map[string]string)

	ctx["Name"] = "Go"

	t, _ := template.ParseFiles("template/index.html")

	err := t.Execute(w, ctx)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User Info from API Server"))
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := r.Context().Value("post")

	t, _ := template.ParseFiles("template/pages/post_detail.html")

	ctx := make(map[string]interface{})
	ctx["post"] = post
	err := t.Execute(w, ctx)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var posts []model.Post

	sql := "SELECT * from posts"

	rows, err := database.DBConn.Query(sql)
	if err != nil {
		log.Println("Error executing query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		data := model.Post{}

		err = rows.Scan(&data.Id, &data.Title, &data.Description)
		if err != nil {
			log.Println("Error scanning rows:", err)
		}
		posts = append(posts, data)
	}

	ctx := make(map[string]interface{})

	ctx["posts"] = posts
	ctx["heading"] = "Article List"

	t, _ := template.ParseFiles("template/pages/post.html")

	err = t.Execute(w, ctx)
	if err != nil {
		log.Println("Error executing template:", err)
	}

}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := make(map[string]interface{})

	// post part
	if r.Method == "POST" {
		r.ParseForm()
		title := r.PostForm.Get("title")
		description := r.PostForm.Get("description")

		stmt := "INSERT INTO posts (title, description) VALUES ($1, $2)"
		q, err := database.DBConn.Prepare(stmt)
		if err != nil {
			log.Println(err)
		}

		res, err := q.Exec(title, description)
		if err != nil {
			log.Println(err)
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 1 {
			ctx["success"] = "Post created successfully!"
		}

		log.Println("Rows affected - ", rowsAffected)
	}

	t, _ := template.ParseFiles("template/pages/post_form.html")

	err := t.Execute(w, ctx)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

func editPostHandler(w http.ResponseWriter, r *http.Request) {

	ctx := make(map[string]interface{})
	p := r.Context().Value("post")

	post := p.(model.Post)

	// Post Part

	if r.Method == "POST" {

		r.ParseForm()

		title := r.PostForm.Get("title")
		description := r.PostForm.Get("description")

		stmt := "UPDATE posts set title=$1, description=$2 where id=$3"

		query, err := database.DBConn.Prepare(stmt)

		if err != nil {
			log.Println(err)
		}

		res, err := query.Exec(title, description, post.Id)

		if err != nil {
			log.Println(err)
		}

		rowsAffected, _ := res.RowsAffected()

		if rowsAffected == 1 {
			ctx["success"] = "Post successully updated."
		}

		log.Println(rowsAffected)

	}

	// Get Part

	// Load template
	t, _ := template.ParseFiles("templates/pages/post_form.html")

	ctx["post"] = post
	err := t.Execute(w, ctx)

	if err != nil {
		log.Println("Error in tpl execution", err)
	}

}
