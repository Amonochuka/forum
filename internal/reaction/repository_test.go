package reaction

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

	schema, err := os.ReadFile("../../migrations/tables.sql")
	if err != nil {
		t.Fatalf("Failed to read schema: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		t.Fatalf("Failed to execute schema: %v", err)
	}

	// Seed user
	_, err = db.Exec(`
		INSERT INTO users (id, email, username, password_hash)
		VALUES (1, 'test@example.com', 'testuser', 'hash')
	`)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Seed post
	_, err = db.Exec(`
		INSERT INTO posts (id, user_id, title, content)
		VALUES (1, 1, 'Test Post', 'Content')
	`)
	if err != nil {
		t.Fatalf("Failed to insert post: %v", err)
	}

	// Seed comment
	_, err = db.Exec(`
		INSERT INTO comments (id, user_id, post_id, content)
		VALUES (1, 1, 1, 'Test Comment')
	`)
	if err != nil {
		t.Fatalf("Failed to insert comment: %v", err)
	}

	return db
}

func TestAddAndGetPostReaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	postID := 1
	reaction := &Reaction{
		UserID: 1,
		PostID: &postID,
		Type:   1,
	}

	err := repo.AddReaction(reaction)
	if err != nil {
		t.Fatalf("AddReaction failed: %v", err)
	}

	r, err := repo.GetUserReaction(1, &postID, nil)
	if err != nil {
		t.Fatalf("GetUserReaction failed: %v", err)
	}

	if r == nil {
		t.Fatalf("expected reaction, got nil")
	}

	if r.Type != 1 {
		t.Errorf("expected reaction type 1, got %d", r.Type)
	}
}

func TestUpdateReaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	postID := 1

	reaction := &Reaction{
		UserID: 1,
		PostID: &postID,
		Type:   1,
	}

	_ = repo.AddReaction(reaction)

	reaction.Type = -1
	err := repo.UpdateReaction(reaction)
	if err != nil {
		t.Fatalf("UpdateReaction failed: %v", err)
	}

	r, _ := repo.GetUserReaction(1, &postID, nil)
	if r.Type != -1 {
		t.Errorf("expected updated reaction -1, got %d", r.Type)
	}
}

func TestDeleteReaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	postID := 1

	reaction := &Reaction{
		UserID: 1,
		PostID: &postID,
		Type:   1,
	}

	_ = repo.AddReaction(reaction)

	err := repo.DeleteReaction(1, &postID, nil)
	if err != nil {
		t.Fatalf("DeleteReaction failed: %v", err)
	}

	r, err := repo.GetUserReaction(1, &postID, nil)
	if err != nil {
		t.Fatalf("GetUserReaction failed: %v", err)
	}
	if r != nil {
		t.Fatalf("expected reaction to be deleted")
	}
}

func TestGetPostReactionCounts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	postID := 1

	_ = repo.AddReaction(&Reaction{UserID: 1, PostID: &postID, Type: 1})

	likes, dislikes, err := repo.GetPostReactionCounts(postID)
	if err != nil {
		t.Fatalf("GetPostReactionCounts failed: %v", err)
	}

	if likes != 1 || dislikes != 0 {
		t.Fatalf("expected 1 like, 0 dislikes — got %d likes, %d dislikes", likes, dislikes)
	}
}

func TestAddAndGetCommentReaction(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	commentID := 1

	reaction := &Reaction{
		UserID:    1,
		CommentID: &commentID,
		Type:      1,
	}

	err := repo.AddReaction(reaction)
	if err != nil {
		t.Fatalf("AddReaction failed: %v", err)
	}

	reactions, err := repo.GetCommentReactions(commentID)
	if err != nil {
		t.Fatalf("GetCommentReactions failed: %v", err)
	}

	if len(reactions) != 1 {
		t.Fatalf("expected 1 reaction, got %d", len(reactions))
	}
}

func TestDeleteReactionNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	postID := 1
	err := repo.DeleteReaction(1, &postID, nil)
	if err == nil {
		t.Fatalf("expected error when deleting missing reaction")
	}
}