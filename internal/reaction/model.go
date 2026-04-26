package reaction

import "time"

type Reaction struct {
	ID        int
	UserID    int
	PostID    *int
	CommentID *int
	Type      int//  1 = like, -1 = dislike
	CreatedAt time.Time
}
