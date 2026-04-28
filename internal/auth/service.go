package auth

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo UserRepo
}

var (
	ErrEmailTaken      = errors.New("email already taken")
	ErrUsernameTaken   = errors.New("username already taken")
	ErrUsernameMissing = errors.New("username is required")
	ErrEmailMissing    = errors.New("email is required")
	ErrPasswordShort   = errors.New("password must be at least 8 characters")
)

func NewService(repo UserRepo) *Service {
	return &Service{Repo: repo}
}

func validateUser(user User) error {
	if strings.TrimSpace(user.Username) == "" {
		return ErrUsernameMissing
	}
	if strings.TrimSpace(user.Email) == "" {
		return ErrEmailMissing
	}
	if len(user.Password) < 8 {
		return ErrPasswordShort
	}
	return nil
}

func (s *Service) Register(user User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	_, err := s.Repo.GetUserByEmail(user.Email)
	if err == nil {
		fmt.Println("Email exist")
		// no error means user was found — email already taken
		return ErrEmailTaken
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
	if err := s.Repo.CreateUser(user); err != nil {
		switch {
		case errors.Is(err, ErrEmailExists):
			return ErrEmailTaken
		case errors.Is(err, ErrUsernameExists):
			return ErrUsernameTaken
		default:
			return err
		}
	}
	return nil
}

func (s *Service) Login(email, password string) (User, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		// generic message — don't reveal whether email exists
		return User{}, ErrInvalidPassword
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("Bcrypt Error:", err)
		return User{}, ErrInvalidPassword
	}

	return user, nil
}

// auth/service.go

func (s *Service) GetUserByID(id int) (User, error) {
	// Simply calls the repository
	return s.Repo.GetUserByID(id)
}
