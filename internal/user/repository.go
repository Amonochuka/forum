package user

import (
	"database/sql"
	"errors"
	"log"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrInternal   = errors.New("internal error")
)

type Repository interface {
	GetByID(id int) (*Profile, error)
}

type sqliteRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepo{db: db}
}

func (r *sqliteRepo) GetByID(id int) (*Profile, error) {
	query := `SELECT id, username FROM users WHERE id = ?`

	user := &Profile{}
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		log.Println("unexpected database error:", err)
		return nil, ErrInternal
	}
	return user, nil
}