package comment

import "net/http"

func RegisterRoutes(h *Handler, requireAuth func(http.Handler) http.Handler) {
	http.Handle("GET /posts/{id}/comments", http.HandlerFunc(h.GetComments))
	http.Handle("POST /posts/{id}/comments", requireAuth(http.HandlerFunc(h.CreateComment)))
	http.Handle("GET /comments/{id}/replies", http.HandlerFunc(h.GetReplies))
	http.Handle("POST /comments/{id}/replies", requireAuth(http.HandlerFunc(h.CreateReply)))
}
