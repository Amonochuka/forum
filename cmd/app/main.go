package main

import (
	"database/sql"
	"log"
	"net/http"

	"forum/internal/auth"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("STARTING SERVER...")
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	repo := auth.NewRepository(db)
	service := auth.NewService(repo)
	handler := auth.NewHandler(service)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	auth.RegisterRoutes(handler)

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
