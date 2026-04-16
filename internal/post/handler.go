package post

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	service *PostService // getting logged-in users
}


//Routes inside the Handler it get to decide which action to take based on the http method. 
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
//returns a list of posts
func (handler *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category") //filter by category
	user := r.URL.Query().Get("user") // filter by user

	posts, err := handler.service.GetAllPosts(category, user)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type","application/json")
	json.NewEncoder(w).Encode(posts)
}

//return a single post by that specific id
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
	panic("unimplemented")
}

func (handler *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
