package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	categories := []string{"Philosophy", "Science", "Art", "Technology"}

	for _, name := range categories {
		_, err := db.Exec("INSERT OR IGNORE INTO categories (name) VALUES (?)", name)
		if err != nil {
			log.Printf("Failed to insert %s: %v", name, err)
		} else {
			log.Printf("Inserted category: %s", name)
		}
	}

	log.Println("Database seeded successfully!")
}
