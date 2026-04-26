package comment

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"forum/internal/session"
	"forum/internal/shared/middleware"
)

type Handler struct {
	service Service
}

func NewHandler(service Service, sessionService *session.Service) *Handler {
	return &Handler{service: service}
}

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

func toCommentView(c Comment) CommentView {
	return CommentView{
		ID:         c.ID,
		AuthorName: c.Name,
		Body:       c.Content,
		Likes:      c.Likes,
		Dislikes:   c.Dislikes,
		CreatedAt:  formatTime(c.CreatedAt),
		ReplyCount: c.ReplyCount,
	}
}

// GET /posts/{id}/comments
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			page = parsed
		}
	}

	comments, count, err := h.service.GetTopLevelComments(postID, page)
	if err != nil {
		http.Error(w, "could not fetch comments", http.StatusInternalServerError)
		return
	}

	views := make([]CommentView, 0, len(comments))

	for _, c := range comments {
		views = append(views, toCommentView(c))
	}

	resp := struct {
		Comments []CommentView `json:"comments"`
		Total    int           `json:"total"`
	}{
		Comments: views,
		Total:    count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// POST /posts/{id}/comments
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")

	comment, err := h.service.CreateComment(userID, postID, content, nil)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyContent):
			http.Error(w, "comment cannot be empty", http.StatusBadRequest)
		case errors.Is(err, ErrContentTooLong):
			http.Error(w, "comment too long", http.StatusBadRequest)
		default:
			http.Error(w, "could not create comment", http.StatusInternalServerError)
		}
		return
	}

	view := toCommentView(*comment)

	resp := struct{
		CommentView CommentView `json:"comment"`
	}{
		CommentView: view,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GET /comments/{id}/replies
func (h *Handler) GetReplies(w http.ResponseWriter, r *http.Request) {
	commentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	replies, err := h.service.GetReplies(commentID)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidParentID):
			http.Error(w, "invalid comment id", http.StatusBadRequest)
		default:
			http.Error(w, "could not fetch replies", http.StatusInternalServerError)
		}
		return
	}

	views := make([]CommentView, 0, len(replies))

	for _, r := range replies {
		views = append(views, toCommentView(r))
	}

	resp := struct {
		Replies []CommentView `json:"replies"`
	}{
		Replies: views,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// POST /comments/{id}/replies
func (h *Handler) CreateReply(w http.ResponseWriter, r *http.Request) {
	parentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	comment, err := h.service.CreateComment(userID, postID, content, &parentID)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyContent):
			http.Error(w, "reply cannot be empty", http.StatusBadRequest)
		case errors.Is(err, ErrContentTooLong):
			http.Error(w, "reply too long", http.StatusBadRequest)
		case errors.Is(err, ErrNestedReplyNotAllowed):
			http.Error(w, "cannot reply to a reply", http.StatusBadRequest)
		case errors.Is(err, ErrInvalidParentID):
			http.Error(w, "invalid parent comment", http.StatusBadRequest)
		default:
			http.Error(w, "could not create reply", http.StatusInternalServerError)
		}
		return
	}

	view := toCommentView(*comment)
	
	resp := struct{
		CommentView CommentView `json:"comment"`
	}{
		CommentView: view,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
