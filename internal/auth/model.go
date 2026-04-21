package auth

import "time"

type User struct {
	ID int
	Username string
	Email string
	Password string //hashed
	CreatedAt time.Time
}
