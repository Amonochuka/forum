package session

import (
	"database/sql"
	"os"
	"testing"
	"time"

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

	// Insert test user (sessions.user_id FK)
	_, err = db.Exec(`
		INSERT INTO users (id, email, username, password_hash)
		VALUES (1, 'test@example.com', 'testuser', 'hash')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	return db
}

func TestCreateAndGetSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	uuid := "session-uuid-123"
	userID := 1
	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour)

	err := repo.CreateSessionRepository(uuid, userID, createdAt, expiresAt)
	if err != nil {
		t.Fatalf("CreateSessionRepository failed: %v", err)
	}

	gotUserID, gotExpiresAt, err := repo.Get(uuid)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if gotUserID != userID {
		t.Errorf("expected userID %d, got %d", userID, gotUserID)
	}

	if !gotExpiresAt.Equal(expiresAt) {
		t.Errorf("expected expiresAt %v, got %v", expiresAt, gotExpiresAt)
	}
}

func TestDeleteSession(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	uuid := "session-to-delete"
	userID := 1
	now := time.Now()

	err := repo.CreateSessionRepository(uuid, userID, now, now.Add(time.Hour))
	if err != nil {
		t.Fatalf("CreateSessionRepository failed: %v", err)
	}

	err = repo.Delete(uuid)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, _, err = repo.Get(uuid)
	if err == nil {
		t.Fatalf("expected error when getting deleted session")
	}
}

func TestDeleteAllUserSessions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)

	userID := 1
	now := time.Now()

	_, _ = repo.db.Exec(`
		INSERT INTO sessions (id, user_id, created_at, expires_at)
		VALUES
		('s1', ?, ?, ?),
		('s2', ?, ?, ?)
	`, userID, now, now.Add(time.Hour), userID, now, now.Add(time.Hour))

	err := repo.DeleteAllUserSessions(userID)
	if err != nil {
		t.Fatalf("DeleteAllUserSessions failed: %v", err)
	}

	rows, err := db.Query(`SELECT COUNT(*) FROM sessions WHERE user_id = ?`, userID)
	if err != nil {
		t.Fatalf("count query failed: %v", err)
	}
	defer rows.Close()

	var count int
	rows.Next()
	rows.Scan(&count)

	if count != 0 {
		t.Fatalf("expected 0 sessions, got %d", count)
	}
}