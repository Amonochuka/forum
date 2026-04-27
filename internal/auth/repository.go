package auth

import (
	"database/sql"
	"errors"
)
type UserRepo interface {
    GetUserByEmail(email string) (User, error)
    CreateUser(user User) error
    GetUserByID(id int) (User, error)
}
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidPassword = errors.New("invalid email or password")
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(user User) error {
	_, err := r.DB.Exec(
		"INSERT INTO users(username, email, password_hash) VALUES (?, ?, ?)",
		user.Username,
		user.Email,
		user.Password,
	)
	return err
}

func (r *Repository) GetUserByEmail(email string) (User, error) {
	var user User

	err := r.DB.QueryRow(
		"SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?",
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	return user, err
}
// auth/repository.go

func (r *Repository) GetUserByID(id int) (User, error) {
    var user User
    // Adjust the query to match your actual database column names
    query := "SELECT id, username, email, password_hash FROM users WHERE id = ?"
    
    err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err != nil {
        return User{}, err // Returns error if not found or DB failure
    }
    return user, nil
}