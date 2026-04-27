package post

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test DB: %v", err)
	}

	// Read and execute migrations
	schema, err := os.ReadFile("../../migrations/tables.sql")
	if err != nil {
		t.Fatalf("Failed to read schema: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatalf("Failed to execute schema: %v", err)
	}

	// Insert test user to satisfy foreign key constraints
	_, err = db.Exec(`INSERT INTO users (id, email, username, password_hash) VALUES (1, 'test@example.com', 'testuser', 'hash')`)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	return db
}

func TestCreateAndGetPost(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPostRepository(db)

	// Test CreatePost
	postID, err := repo.CreatePost(1, "Test Title", "Test Content")
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if postID == 0 {
		t.Fatalf("Expected valid postID, got 0")
	}

	// Test GetPostByID
	post, err := repo.GetPostByID(postID)
	if err != nil {
		t.Fatalf("GetPostByID failed: %v", err)
	}

	if post.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got '%s'", post.Title)
	}

	if post.Content != "Test Content" {
		t.Errorf("Expected content 'Test Content', got '%s'", post.Content)
	}

	if post.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", post.UserID)
	}
}

func TestGetPostByUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPostRepository(db)

	_, err := repo.CreatePost(1, "User Post 1", "Content 1")
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	_, err = repo.CreatePost(1, "User Post 2", "Content 2")
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	posts, err := repo.GetPostByUser(1)
	if err != nil {
		t.Fatalf("GetPostByUser failed: %v", err)
	}

	if len(posts) != 2 {
		t.Fatalf("Expected 2 posts, got %d", len(posts))
	}
}
