package comment

import (
	"time"

	"forum/internal/auth"
)

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

// View model - how the data is presented to the user
type CommentsSectionData struct {
	PostID      int
	CurrentUser auth.User
	Comments    []Comment
	TotalCount  int
}
