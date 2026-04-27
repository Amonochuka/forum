package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/comment"
	"forum/internal/post"
	"forum/internal/reaction"
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
	authHandler := auth.NewHandler(authService, sessionService)
	auth.RegisterRoutes(authHandler)

	requireAuth := middleware.RequireAuth(sessionService)
	optionalAuth := middleware.OptionalAuth(sessionService)

	// comments
	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService, sessionService)
	comment.RegisterRoutes(commentHandler, requireAuth)
	//post
	postRepo := post.NewPostRepository(db)
	catRepo := post.NewCategoryRepository(db)
	userRepo := post.NewUserRepository(db)
	reactionRepo := reaction.NewRepository(db)

	postservice := post.NewPostService(postRepo, catRepo, userRepo, reactionRepo)
	posthandler := post.NewPostHandler(postservice)
	post.RegisterPostRoutes(posthandler, requireAuth, optionalAuth)

	// reaction
	reactionService := &reaction.ReactionService{Repo: reactionRepo}
	reactionHandler := reaction.NewHandler(reactionService)
	reaction.RegisterRoutes(reactionHandler, requireAuth)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
