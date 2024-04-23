package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/user-info", userInfoHandler)

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
