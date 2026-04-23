package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// helper to talsk to the repository
type Service struct {
	repo *Repository
}

// The Constructor: This connects the Brain to the Librarian
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// function creates a session when user logins or registers
func (s *Service) StartSession(userid int) (string, error) {
	// generate a unique ID
	token := uuid.New().String()
	created := time.Now()
	// set the expiration time for the token to 24 hrs after creation
	expires := time.Now().Add(24 * time.Hour)
	// save this in the database repository
	err := s.repo.CreateSessionRepository(token, userid, created, expires)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) ValidateSession(token string) (int, error) {
	userID, expiresAt, err := s.repo.Get(token)
	if err != nil {
		return 0, err
	}
	// check if current time is equal to expiry time
	if time.Now().After(expiresAt) {
		s.repo.Delete(token)
		return 0, errors.New("this ticket is expired")
	}
	return userID, nil
}
func (s *Service) DeleteSession(token string) error {
	return s.repo.Delete(token)
}