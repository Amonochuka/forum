package comment

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, requireAuth func(http.HandlerFunc) http.HandlerFunc) {
	mux.HandleFunc("GET /posts/{id}/comments", h.GetComments)
	mux.HandleFunc("POST /posts/{id}/comments", requireAuth(h.CreateComment))
	mux.HandleFunc("GET /comments/{id}/replies", h.GetReplies)
	mux.HandleFunc("POST /comments/{id}/replies", requireAuth(h.CreateReply))
}
