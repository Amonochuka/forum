package comment

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"forum/internal/shared/middleware"
	"forum/internal/user"
)

type Handler struct {
	service     Service
	userService user.Service
}

func NewHandler(service Service, userService user.Service) *Handler {
	return &Handler{service: service, userService: userService}
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

	views := ToCommentViews(comments)

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

	// Get user
	user, err := h.userService.GetByID(userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Add their name
	comment.Name = user.Username

	view := ToCommentView(*comment)

	resp := struct {
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

	views := ToCommentViews(replies)

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

	view := ToCommentView(*comment)

	resp := struct {
		CommentView CommentView `json:"comment"`
	}{
		CommentView: view,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
