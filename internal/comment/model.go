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
	Name            string
	Likes, Dislikes int
	ReplyCount      int
	CreatedAt       time.Time
}

// View model - how the data is presented to the user
type CommentsSectionData struct {
	PostID      int
	CurrentUser auth.User
	Comments    []Comment
	TotalCount  int
}

type CommentView struct {
	ID         int    `json:"id"`
	AuthorName string `json:"authorName"`
	Body       string `json:"body"`
	Likes      int    `json:"likes"`
	Dislikes   int    `json:"dislikes"`
	CreatedAt  string `json:"createdAt"`
	ReplyCount int    `json:"replyCount"`
}
