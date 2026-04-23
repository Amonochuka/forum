package post

import "net/http"

func RegisterPostRoutes(handler *PostHandler, requireAuth func(http.Handler) http.Handler) {
	http.Handle("GET /posts", http.HandlerFunc(handler.GetPosts))
	http.Handle("GET /posts/{id}", http.HandlerFunc(handler.GetPostByID))
	http.Handle("POST /posts", requireAuth(http.HandlerFunc(handler.CreatePost)))
}
