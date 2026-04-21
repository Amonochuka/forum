package comment

import (
	"errors"
	"forum/internal/auth"
	"net/http"
	"strconv"
)

// TODO: Replace with actual session service
type sessionService interface {
	GetUserFromRequest(r *http.Request) (*auth.User, error)
}

type Handler struct {
	service        Service
	sessionService sessionService
}

func NewHandler(service Service, sessionService sessionService) *Handler {
	return &Handler{service: service, sessionService: sessionService}
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

	comments, err := h.service.GetTopLevelComments(postID, page)
	if err != nil {
		http.Error(w, "could not fetch comments", http.StatusInternalServerError)
		return
	}

	// TODO: Render template with comments
	_ = comments
}

// POST /posts/{id}/comments
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	// TODO: Swap this out when session service is available
	user, err := h.sessionService.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")

	if err := h.service.CreateComment(user.ID, postID, content, nil); err != nil {
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

	http.Redirect(w, r, "/posts/"+strconv.Itoa(postID)+"/comments", http.StatusSeeOther)
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

	// TODO: Render template with replies
	_ = replies
}

// POST /comments/{id}/replies
func (h *Handler) CreateReply(w http.ResponseWriter, r *http.Request) {
	parentID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	// TODO: Swap this out when session service is available
	user, err := h.sessionService.GetUserFromRequest(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateComment(user.ID, postID, content, &parentID); err != nil {
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

	http.Redirect(w, r, "/posts/"+strconv.Itoa(postID)+"/comments", http.StatusSeeOther)
}
