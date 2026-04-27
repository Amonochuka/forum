package reaction

import "time"

type Reaction struct {
	ID        int  `json:"id"`
	UserID    int `json:"user_id`
	PostID    *int `json:"post_id"`
	CommentID *int `json:"comment_id`
	Type      int `json:"reaction_type"`//  1 = like, -1 = dislike
	CreatedAt time.Time
}
