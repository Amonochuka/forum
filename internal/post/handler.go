package post

import (
	"encoding/json"
	"fmt"
	"forum/internal/comment"
	"forum/internal/shared/middleware"
	"html/template"
	"net/http"
	"strconv"
)

type PostHandler struct {
	Service        *PostService
	commentService comment.Service
	templates      *template.Template
}

type PostDetailPageData struct {
	Post       PostResponse
	PostID     int
	Comments   []comment.CommentView
	TotalCount int
}

type CreatePostRequest struct {
	Title    string   `json:"Title"`
	Content  string   `json:"Content"`
	Category []string `json:"category"`
}

func NewPostHandler(service *PostService, commentService comment.Service, templates *template.Template) *PostHandler {
	return &PostHandler{
		Service:        service,
		commentService: commentService,
		templates:      templates}
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
		fmt.Println("Error:", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
		return
	}
}

// return a single post page by that specific id
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

	comments, total, err := handler.commentService.GetTopLevelComments(id, 1)
	if err != nil {
		http.Error(w, "could not fetch comments", http.StatusInternalServerError)
		return
	}

	pageData := PostDetailPageData{
		Post:       post,
		PostID:     id,
		Comments:   comment.ToCommentViews(comments),
		TotalCount: total,
	}

	if err := handler.templates.ExecuteTemplate(w, "post_detail", pageData); err != nil {
		http.Error(w, "failed to render post detail", http.StatusInternalServerError)
	}
}

// return a single post as JSON
func (handler *PostHandler) GetPostByIDAPI(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "JSON encode error", http.StatusInternalServerError)
	}
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
