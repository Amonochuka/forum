package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func Init(dbPath, schemaPath, seedPath string) (*sql.DB, error) {
	// ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}

	// check if DB exists
	_, err := os.Stat(dbPath)
	dbExists := err == nil

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if !dbExists {
		log.Println("DB not found, initializing")

		schema, err := os.ReadFile(schemaPath)
		if err != nil {
			return nil, err
		}
		if _, err := db.Exec(string(schema)); err != nil {
			return nil, err
		}

		seed, err := os.ReadFile(seedPath)
		if err != nil {
			return nil, err
		}
		if _, err := db.Exec(string(seed)); err != nil {
			return nil, err
		}
	}

	return db, nil
}