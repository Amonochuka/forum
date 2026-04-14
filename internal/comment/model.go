package comment

import "time"

type Comment struct {
	ID int
	UserID int
	PostID int
	Content string
	name string
	likes, dislikes int
	CreatedAt time.Time
}