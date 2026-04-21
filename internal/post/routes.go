package post

import "net/http"

func RegisterPostRoutes(handler *PostHandler) {
	http.HandleFunc("/posts", handler.HandlePosts)
	http.HandleFunc("/posts/", handler.GetPostByID)
	
}
