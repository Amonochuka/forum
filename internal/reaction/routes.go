package reaction

import "net/http"

func RegisterRoutes(handler *Handler, requireAuth func(http.Handler) http.Handler) {
	http.Handle("/react", requireAuth(http.HandlerFunc(handler.React)))
	http.Handle("/comments/{id}/reactions", http.HandlerFunc(handler.GetCommentReactionCounts))
}