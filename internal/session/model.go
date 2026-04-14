package session

import "time"

type session struct{
	ID int
	UserID int
	CreatedAt time.Time
	ExpiresAt time.Time
}