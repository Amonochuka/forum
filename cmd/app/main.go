package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/comment"
	"forum/internal/session"
	"forum/internal/shared/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	auth.RegisterRoutes(authHandler)

	// session
	sessionRepo := session.NewRepository(db)
	sessionService := session.NewService(sessionRepo)

	requireAuth := middleware.RequireAuth(sessionService)

	// comments
	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService, sessionService)
	comment.RegisterRoutes(commentHandler, requireAuth)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
