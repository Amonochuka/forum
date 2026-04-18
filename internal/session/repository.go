package session

import (
	"database/sql"
	"time"
)

// Repository is a "box" that holds our database connection
// acts a bridge to our db and our go code
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSessionRepository(uuid string, userID int, createdAt time.Time, expiresAt time.Time) error {
	query := `INSERT INTO sessions(id,user_id,created_at,expires_at) VALUES(?, ?, ?, ?)`
	_, err := r.db.Exec(query, uuid, userID, createdAt, expiresAt)
	return err
}

func (r *Repository) Get(uuid string) (int, time.Time, error) {
	var userID int
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM sessions WHERE id = ?`

	// 2. QueryRow means "Find exactly one thing"
	err := r.db.QueryRow(query, uuid).Scan(&userID, &expiresAt)

	return userID, expiresAt, err
}

// this deletes the session safely when a user logsout
func (r *Repository) Delete(uuid string) error {
	query := `DELETE FROM sessions WHERE id = ?`

	_, err := r.db.Exec(query, uuid)
	return err
}
