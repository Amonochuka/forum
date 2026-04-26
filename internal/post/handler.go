package post

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"forum/internal/shared/middleware"
)

type PostHandler struct {
	Service *PostService
}

type CreatePostRequest struct {
	Title    string   `json:"Title"`
	Content  string   `json:"Content"`
	Category []string `json:"category"`
}

func NewPostHandler(service *PostService) *PostHandler {
	return &PostHandler{Service: service}
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
	category := r.URL.Query().Get("category")
	user := r.URL.Query().Get("user")
	liked := r.URL.Query().Get("liked")

	var likedBy string
	if liked == "true" {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		likedBy = strconv.Itoa(userID)
	}

	// Also support "user=me" for created posts
	if user == "me" {
		userID, ok := middleware.GetUserID(r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user = strconv.Itoa(userID)
	}

	posts, err := handler.Service.GetPosts(category, user, likedBy)
	if err != nil {
		fmt.Println("Error:",err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
		return
	}
}

// return a single post by that specific id
func (handler *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	post, err := handler.Service.GetPostByID(id)
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

	err = handler.Service.CreatePost(userID, req.Title, req.Content, req.Category)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *PostHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := handler.Service.GetAllCategories()
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
