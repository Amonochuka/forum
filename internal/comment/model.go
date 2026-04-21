package comment

import "time"

type Comment struct {
	ID              int
	UserID          int
	PostID          int
	ParentID        *int
	Content         string
	name            string
	Likes, Dislikes int
	CreatedAt       time.Time
}
