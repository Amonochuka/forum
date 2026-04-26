package auth

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func validateUser(user User) error {
	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username is required")
	}
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func (s *Service) Register(user User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	_, err := s.Repo.GetUserByEmail(user.Email)
	if err == nil {
		// no error means user was found — email already taken
		return ErrEmailExists
	}
	if !errors.Is(err, ErrUserNotFound) {
		// real DB error — don't proceed
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.Repo.CreateUser(user)
}

func (s *Service) Login(email, password string) (User, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		// generic message — don't reveal whether email exists
		return User{}, ErrInvalidPassword
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, ErrInvalidPassword
	}

	return user, nil
}
