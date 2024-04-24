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

		// r.Get("/", getPostHandler)
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

	t, _ := template.ParseFiles("template/post_from.html")

	err := t.Execute(w, ctx)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
