package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"forum/internal/auth"
	"forum/internal/comment"
	"forum/internal/post"
	"forum/internal/reaction"
	"forum/internal/session"
	"forum/internal/shared/middleware"
	"forum/internal/user"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/app.db")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles(
		"web/templates/index.html",
		"web/templates/post_feed.html",
		"web/templates/post_detail.html",
		"web/templates/components/navbar.html",
		"web/templates/components/hero.html",
		"web/templates/components/create_post.html",
		"web/templates/components/sidebar.html",
		"web/templates/components/footer.html",
		"web/templates/components/scripts.html",
		"web/templates/components/comments_section.html",
		"web/templates/components/comment.html",
	)
	if err != nil {
		log.Fatal("failed to parse templates:", err)
	}

	// auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// session
	sessionRepo := session.NewRepository(db)
	sessionService := session.NewService(sessionRepo)
	authHandler := auth.NewHandler(authService, sessionService, tmpl)
	auth.RegisterRoutes(authHandler)

	requireAuth := middleware.RequireAuth(sessionService)

	// user
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// comments
	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService, userService)
	comment.RegisterRoutes(commentHandler, requireAuth)
	//post
	postRepo := post.NewPostRepository(db)
	catRepo := post.NewCategoryRepository(db)
	userRepo := post.NewUserRepository(db)
	reactionRepo := reaction.NewRepository(db)

	postservice := post.NewPostService(postRepo, catRepo, userRepo, reactionRepo)
	posthandler := post.NewPostHandler(postservice, commentService, tmpl)
	post.RegisterPostRoutes(posthandler, requireAuth)

	// reaction
	reactionService := &reaction.ReactionService{Repo: reactionRepo}
	reactionHandler := reaction.NewHandler(reactionService)
	reaction.RegisterRoutes(reactionHandler, requireAuth)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
