package post

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/internal/comment"
	"forum/internal/reaction"
	"forum/internal/shared/middleware"
)

func TestPostHandler_CreatePost(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tmpl := template.Must(template.New("test").Parse(`hello`))

	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	postRepo := NewPostRepository(db)
	catRepo := NewCategoryRepository(db)
	userRepo := NewUserRepository(db)
	reactionRepo := &reaction.ReactionRepository{DB: db}
	service := NewPostService(postRepo, catRepo, userRepo, reactionRepo, commentRepo)
	handler := NewPostHandler(service, commentService, tmpl)

	db.Exec(`INSERT INTO categories (id, name) VALUES (1, 'tech')`)

	reqBody := CreatePostRequest{
		Title:    "Test API Title",
		Content:  "Test API Content",
		Category: []string{"tech"},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(bodyBytes))
	
	// Add user context (mocking middleware)
	req = middleware.WithUserID(req, 1)

	rr := httptest.NewRecorder()

	handler.CreatePost(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Verify post was created
	posts, _ := service.GetPosts("", "", "")
	if len(posts) != 1 {
		t.Fatalf("Expected 1 post to be created, got %d", len(posts))
	}
	if posts[0].Title != "Test API Title" {
		t.Errorf("Expected title 'Test API Title', got '%s'", posts[0].Title)
	}
}

func TestPostHandler_GetPosts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tmpl := template.Must(template.New("test").Parse(`hello`))
	
	commentRepo := comment.NewRepository(db)
	commentService := comment.NewService(commentRepo)
	postRepo := NewPostRepository(db)
	catRepo := NewCategoryRepository(db)
	userRepo := NewUserRepository(db)
	reactionRepo := &reaction.ReactionRepository{DB: db}
	service := NewPostService(postRepo, catRepo, userRepo, reactionRepo, commentRepo)
	handler := NewPostHandler(service, commentService, tmpl)

	db.Exec(`INSERT INTO categories (id, name) VALUES (1, 'tech')`)
	service.CreatePost(1, "Post 1", "Content", []string{"tech"})
	service.CreatePost(1, "Post 2", "Content", []string{"tech"})

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rr := httptest.NewRecorder()

	handler.GetPosts(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []PostResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 posts in response, got %d", len(response))
	}
}
