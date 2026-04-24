package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/comment"
	"forum/internal/post"
	"forum/internal/session"
	"forum/internal/shared/middleware"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/app.db")
	if err != nil {
		log.Fatal(err)
	}

	// auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// session
	sessionRepo := session.NewRepository(db)
	sessionService := session.NewService(sessionRepo)
	authHandler := auth.NewHandler(authService,sessionService)
	auth.RegisterRoutes(authHandler)

	requireAuth := middleware.RequireAuth(sessionService)

	// comments
	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService, sessionService)
	comment.RegisterRoutes(commentHandler, requireAuth)
//post

	postRepo:= post.NewPostRepository(db)
	userRepo:= post.NewUserRepository(db)
	catRepo:= post.NewCategoryRepository(db)

	postservice := post.NewPostService(postRepo,userRepo,catRepo)
	posthandler := post.NewPostHandler(postservice)
	post.RegisterPostRoutes(posthandler, requireAuth)
	
	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
