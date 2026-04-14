package post 

import "time"

type Post struct {
	ID int
	UserID int
	PostID int
	Content string
	name string
	categories []Category
	likes, dislikes int
	CreatedAt time.Time

}

type Category struct {
	ID int
	Name string
}
