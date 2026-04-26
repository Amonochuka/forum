package reaction

import (
	"forum/internal/shared/middleware"
	"net/http"
	"strconv"
)

// Handler holds dependencies
type Handler struct {
	ReactionService *ReactionService
}

func NewHandler(service *ReactionService) *Handler {
	return &Handler{ReactionService: service}
}

// React handles like/dislike actions
func (h *Handler) React(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get user from context (set by middleware)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	reactionType := r.FormValue("type")

	var postID *int
	var commentID *int

	// parse post_id
	if pid := r.FormValue("post_id"); pid != "" {
		id, err := strconv.Atoi(pid)
		if err == nil {
			postID = &id
		}
	}

	// parse comment_id
	if cid := r.FormValue("comment_id"); cid != "" {
		id, err := strconv.Atoi(cid)
		if err == nil {
			commentID = &id
		}
	}

	// build reaction model
	reaction := &Reaction{
		UserID:    userID,
		PostID:    postID,
		CommentID: commentID,
		Type:      reactionType,
	}

	// call service
	err := h.ReactionService.React(reaction)
	if err != nil {
		http.Error(w, "failed to react", http.StatusInternalServerError)
		return
	}

	// redirect back to previous page
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}
