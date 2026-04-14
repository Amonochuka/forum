package reaction

type Reaction struct{
	ID int
	UserID int
	PostID int
	CommentID *int
	Type string // "like" or "dislike"

}