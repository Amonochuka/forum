package post

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"forum/internal/shared/middleware"
)

type PostHandler struct {
	service        *PostService 
}	

type CreatePostRequest struct {
	Title string `json:"Title"`
	Content string `json:"Content"`
	Category []string `json:"category"`

}


// Routes inside the Handler it get to decide which action to take based on the http method.
func (handler *PostHandler) HandlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handler.GetPosts(w, r)

	case http.MethodPost:
		handler.CreatePost(w, r)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// returns a list of posts
func (handler *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category") // filter by category
	user := r.URL.Query().Get("user")         // filter by user

	posts, err := handler.service.GetPosts(category, user)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// return a single post by that specific id
func (handler *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/posts/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	post, err := handler.service.GetPostByID(id)
	if err != nil {
		http.Error(w, "Post Not Found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (handler *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = handler.service.CreatePost(userID, req.Title, req.Content, req.Category)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
