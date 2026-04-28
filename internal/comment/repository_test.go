package comment

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test DB: %v", err)
	}

	schema, err := os.ReadFile("../../migrations/tables.sql")
	if err != nil {
		t.Fatalf("failed to read schema: %v", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("failed to execute schema: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO users (id, email, username, password_hash) VALUES (1, 'one@example.com', 'user1', 'hash')`); err != nil {
		t.Fatalf("failed to insert first user: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO users (id, email, username, password_hash) VALUES (2, 'two@example.com', 'user2', 'hash')`); err != nil {
		t.Fatalf("failed to insert second user: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO posts (id, user_id, title, content) VALUES (1, 1, 'Post title', 'Post content')`); err != nil {
		t.Fatalf("failed to insert post: %v", err)
	}

	return db
}

func TestRepository_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	comment := &Comment{
		UserID:  1,
		PostID:  1,
		Content: "First comment",
	}

	if err := repo.Create(comment); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if comment.ID == 0 {
		t.Fatal("expected created comment ID to be set")
	}

	got, err := repo.GetByID(comment.ID)
	if err != nil {
		t.Fatalf("get by ID failed: %v", err)
	}

	if got.Content != "First comment" {
		t.Fatalf("expected content %q, got %q", "First comment", got.Content)
	}
	if got.ParentID != nil {
		t.Fatalf("expected top-level comment to have nil parent, got %v", *got.ParentID)
	}
}

func TestRepository_GetTopLevelByPostWithReactions_IncludesCounts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	if _, err := db.Exec(`INSERT INTO comments (id, post_id, user_id, content) VALUES (10, 1, 1, 'Top level')`); err != nil {
		t.Fatalf("failed to insert top-level comment: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO comments (id, post_id, user_id, parent_id, content) VALUES (11, 1, 2, 10, 'Reply')`); err != nil {
		t.Fatalf("failed to insert reply: %v", err)
	}
	if _, err := db.Exec(`INSERT INTO reactions (user_id, comment_id, reaction_type) VALUES (1, 10, 1), (2, 10, -1)`); err != nil {
		t.Fatalf("failed to insert reactions: %v", err)
	}

	comments, err := repo.GetTopLevelByPostWithReactions(1, 20, 0)
	if err != nil {
		t.Fatalf("get top-level comments failed: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 top-level comment, got %d", len(comments))
	}

	got := comments[0]
	if got.Name != "user1" {
		t.Fatalf("expected username %q, got %q", "user1", got.Name)
	}
	if got.Likes != 1 || got.Dislikes != 1 {
		t.Fatalf("expected 1 like and 1 dislike, got %d likes and %d dislikes", got.Likes, got.Dislikes)
	}
	if got.ReplyCount != 1 {
		t.Fatalf("expected reply count 1, got %d", got.ReplyCount)
	}
}
