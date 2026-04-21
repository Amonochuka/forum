package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	root := "."

	schemaPath := filepath.Join(root, "migrations", "tables.sql")
	seedPath := filepath.Join(root, "internal", "db", "seed.sql")
	dbDir := filepath.Join(root, "data")
	dbPath := filepath.Join(dbDir, "app.db")

	// 1. ensure /data exists
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatal(err)
	}

	// 2. remove ONLY the db file (not the whole folder)
	if err := os.Remove(dbPath); err == nil {
		fmt.Println("old db removed")
	} else if !os.IsNotExist(err) {
		log.Fatal(err)
	}

	// (optional but smart) remove WAL/SHM too if they exist
	os.Remove(dbPath + "-wal")
	os.Remove(dbPath + "-shm")

	// 3. open db (creates new one)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 4. run schema
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatal(err)
	}

	// 5. run seed
	seed, err := os.ReadFile(seedPath)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(string(seed)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("db initialized successfully")
}
