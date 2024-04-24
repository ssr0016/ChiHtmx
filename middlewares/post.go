package middlewares

import (
	"ChiHtmx/database"
	"ChiHtmx/model"
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var post model.Post

		if id := chi.URLParam(r, "id"); id != "" {
			stmt := "SELECT * FROM posts WHERE id = $1"
			row := database.DBConn.QueryRow(stmt, id)

			err := row.Scan(&post.Id, &post.Title, &post.Description)
			if err != nil {
				log.Println("Error scanning row:", err)
			}
		}

		ctx := context.WithValue(r.Context(), "post", post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
