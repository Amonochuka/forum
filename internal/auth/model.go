package auth

import "time"

type user struct {
	ID int
	Username string
	Email string
	Password string //hashed
	CreatedAt time.Time
}