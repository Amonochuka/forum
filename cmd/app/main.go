package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	// db "forum/cmd/testdb"
	"forum/internal/auth"
	"forum/internal/comment"
	"forum/internal/db"
	"forum/internal/post"
	"forum/internal/reaction"
	"forum/internal/session"
	"forum/internal/shared/middleware"
	"forum/internal/user"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/app.db"
	}

	db, err := db.Init(
		dbPath,
		"./migrations/tables.sql",
		"./internal/db/seed.sql",
	)
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
	optionalAuth := middleware.OptionalAuth(sessionService)

	// user
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// comments
	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	commentHandler := comment.NewHandler(commentService, userService)
	comment.RegisterRoutes(commentHandler, requireAuth)
	// post
	postRepo := post.NewPostRepository(db)
	catRepo := post.NewCategoryRepository(db)
	userRepo := post.NewUserRepository(db)
	reactionRepo := reaction.NewRepository(db)

	postservice := post.NewPostService(postRepo, catRepo, userRepo, reactionRepo)
	posthandler := post.NewPostHandler(postservice, commentService, tmpl)
	post.RegisterPostRoutes(posthandler, requireAuth, optionalAuth)

	// reaction
	reactionService := &reaction.ReactionService{Repo: reactionRepo}
	reactionHandler := reaction.NewHandler(reactionService)
	reaction.RegisterRoutes(reactionHandler, requireAuth)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
