package comment

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler, requireAuth func(http.Handler) http.Handler) {
	mux.Handle("GET /posts/{id}/comments", http.HandlerFunc(h.GetComments))
	mux.Handle("POST /posts/{id}/comments", requireAuth(http.HandlerFunc(h.CreateComment)))
	mux.Handle("GET /comments/{id}/replies", http.HandlerFunc(h.GetReplies))
	mux.Handle("POST /comments/{id}/replies", requireAuth(http.HandlerFunc(h.CreateReply)))
}
