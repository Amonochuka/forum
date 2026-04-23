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

	// repositories
	authRepo := auth.NewRepository(db)
	commentRepo := comment.NewRepository(db)
	sessionRepo := session.NewRepository(db)

	// services
	authService := auth.NewService(authRepo)
	sessionService := session.NewService(sessionRepo)

	// handlers
	authHandler := auth.NewHandler(authService,sessionService)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	auth.RegisterRoutes(authHandler)

	// session

	requireAuth := middleware.RequireAuth(sessionService)

	// comments
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService, sessionService)
	comment.RegisterRoutes(commentHandler, requireAuth)

	//router

	mux:=http.NewServeMux()

	//auth

	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login",authHandler.Login)
	//view-only
	//protected routes
	//startserver
	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
