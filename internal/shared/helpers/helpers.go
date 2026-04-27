package helpers

import (
	"fmt"
	"time"

	"forum/internal/comment"
)

func formatTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"

	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))

	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))

	case diff < 48*time.Hour:
		return "yesterday"

	default:
		return t.Format("02 Jan 2006")
	}
}

func ToCommentView(c comment.Comment) comment.CommentView {
	return comment.CommentView{
		ID:         c.ID,
		AuthorName: c.Name,
		Body:       c.Content,
		Likes:      c.Likes,
		Dislikes:   c.Dislikes,
		CreatedAt:  formatTime(c.CreatedAt),
		ReplyCount: c.ReplyCount,
	}
}
