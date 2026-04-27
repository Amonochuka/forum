package session

import (
	"errors"
	"testing"
	"time"
)

type mockSessionRepo struct {
	createSessionFunc         func(uuid string, userID int, createdAt time.Time, expiresAt time.Time) error
	getFunc                   func(uuid string) (int, time.Time, error)
	deleteFunc                func(uuid string) error
	deleteAllUserSessionsFunc func(userID int) error
}

func (m *mockSessionRepo) CreateSessionRepository(uuid string, userID int, createdAt time.Time, expiresAt time.Time) error {
	return m.createSessionFunc(uuid, userID, createdAt, expiresAt)
}

func (m *mockSessionRepo) Get(uuid string) (int, time.Time, error) {
	return m.getFunc(uuid)
}

func (m *mockSessionRepo) Delete(uuid string) error {
	return m.deleteFunc(uuid)
}

func (m *mockSessionRepo) DeleteAllUserSessions(userID int) error {
	return m.deleteAllUserSessionsFunc(userID)
}

func TestStartSession(t *testing.T) {
	deletedUserID := 0
	createdToken := ""
	createdUserID := 0

	service := NewService(&mockSessionRepo{
		createSessionFunc: func(uuid string, userID int, createdAt time.Time, expiresAt time.Time) error {
			createdToken = uuid
			createdUserID = userID
			if expiresAt.Sub(createdAt) < 23*time.Hour {
				t.Fatalf("expected session expiry about 24h after creation, got %v", expiresAt.Sub(createdAt))
			}
			return nil
		},
		getFunc: func(uuid string) (int, time.Time, error) {
			return 0, time.Time{}, errors.New("not implemented")
		},
		deleteFunc: func(uuid string) error {
			return nil
		},
		deleteAllUserSessionsFunc: func(userID int) error {
			deletedUserID = userID
			return nil
		},
	})

	token, err := service.StartSession(1)
	if err != nil {
		t.Errorf("expected no error,but got:%v", err)
	}
	if token == "" {
		t.Errorf("Expected a token,but got an empty string")
	}
	if deletedUserID != 1 {
		t.Errorf("expected old sessions to be deleted for user 1, got %d", deletedUserID)
	}
	if createdUserID != 1 {
		t.Errorf("expected session to be created for user 1, got %d", createdUserID)
	}
	if createdToken != token {
		t.Errorf("expected created token %q to match returned token %q", createdToken, token)
	}
}

func TestValidateSession(t *testing.T) {
	service := NewService(&mockSessionRepo{
		createSessionFunc: func(uuid string, userID int, createdAt time.Time, expiresAt time.Time) error {
			return nil
		},
		getFunc: func(uuid string) (int, time.Time, error) {
			return 0, time.Time{}, errors.New("session not found")
		},
		deleteFunc: func(uuid string) error {
			return nil
		},
		deleteAllUserSessionsFunc: func(userID int) error {
			return nil
		},
	})

	_, err := service.ValidateSession("1234")
	if err == nil {
		t.Errorf("Expected an error for fake token, but got none")
	}
}
