package post

import (
	"testing"

	"forum/internal/reaction"

	_ "github.com/mattn/go-sqlite3"
)

func TestPostService_CreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Setup repositories
	postRepo := NewPostRepository(db)
	catRepo := NewCategoryRepository(db)
	userRepo := NewUserRepository(db)
	reactionRepo := &reaction.ReactionRepository{DB: db}

	// Setup service
	service := NewPostService(postRepo, catRepo, userRepo, reactionRepo)

	// Create a category
	_, err := db.Exec(`INSERT INTO categories (id, name) VALUES (1, 'tech')`)
	if err != nil {
		t.Fatalf("Failed to insert category: %v", err)
	}

	// Test CreatePost
	err = service.CreatePost(1, "Service Title", "Service Content", []string{"tech"})
	if err != nil {
		t.Fatalf("Failed to create post via service: %v", err)
	}

	// Test GetPostByID
	post, err := service.GetPostByID(1)
	if err != nil {
		t.Fatalf("Failed to get post via service: %v", err)
	}

	if post.Title != "Service Title" {
		t.Errorf("Expected title 'Service Title', got '%s'", post.Title)
	}
	if post.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", post.Username)
	}
	if len(post.Category) != 1 || post.Category[0].Name != "tech" {
		t.Errorf("Expected category 'tech', got %v", post.Category)
	}

	// Test GetPosts (all)
	posts, err := service.GetPosts("", "", "")
	if err != nil {
		t.Fatalf("Failed to get posts: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}
}

func TestPostService_GetPostsByCategory(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	postRepo := NewPostRepository(db)
	catRepo := NewCategoryRepository(db)
	userRepo := NewUserRepository(db)
	reactionRepo := &reaction.ReactionRepository{DB: db}
	service := NewPostService(postRepo, catRepo, userRepo, reactionRepo)

	db.Exec(`INSERT INTO categories (id, name) VALUES (1, 'science')`)
	db.Exec(`INSERT INTO categories (id, name) VALUES (2, 'art')`)

	service.CreatePost(1, "Post 1", "Content 1", []string{"science"})
	service.CreatePost(1, "Post 2", "Content 2", []string{"art"})

	posts, err := service.GetPosts("science", "", "")
	if err != nil {
		t.Fatalf("Failed to get posts by category: %v", err)
	}

	if len(posts) != 1 {
		t.Fatalf("Expected 1 post for science category, got %d", len(posts))
	}
	if posts[0].Title != "Post 1" {
		t.Errorf("Expected Post 1, got %s", posts[0].Title)
	}
}
