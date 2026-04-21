package auth

import (
	"database/sql"
	"errors"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(user User) error {
	_, err := r.DB.Exec(
		"INSERT INTO users(username, email, password) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		user.Password,
	)
	return err
}

func (r *Repository) GetUserByEmail(email string) (User, error) {
	var user User

	err := r.DB.QueryRow(
		"SELECT id, username, email, password, created_at FROM users WHERE email = ?",
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return User{}, errors.New("user not found")
	}

	return user, err
}
