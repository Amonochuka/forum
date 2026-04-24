package post

import "time"

// How the data is store in the database
type Post struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// How the data will be presented to the user (DTO)
type PostResponse struct {
	ID              int
	Title           string
	Content         string
	Username        string
	Category        []Category
	Likes, Dislikes int
	CreatedAt       time.Time
}

type Category struct {
	ID   int
	Name string
}
