package post

import "net/http"

func RegisterPostRoutes(handler *PostHandler, requireAuth func(http.Handler) http.Handler, optionalAuth func(http.Handler) http.Handler) {
	http.Handle("GET /posts", optionalAuth(http.HandlerFunc(handler.GetPosts)))
	http.Handle("GET /posts/{id}", http.HandlerFunc(handler.GetPostByID))
	http.Handle("GET /api/posts/{id}", http.HandlerFunc(handler.GetPostByIDAPI))
	http.Handle("POST /posts", requireAuth(http.HandlerFunc(handler.CreatePost)))
	http.Handle("GET /categories", http.HandlerFunc(handler.GetCategories))
}
