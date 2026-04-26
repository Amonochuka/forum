package auth

import (
	"encoding/json"
	"forum/internal/session"
	"html/template"
	"net/http"
	"time"
)

type Handler struct {
	AuthService *Service
	SessionService *session.Service
	templates *template.Template
}

func NewHandler(authservice *Service,sessionService *session.Service, templates *template.Template) *Handler {
	return &Handler{AuthService: authservice,
		SessionService: sessionService,
		templates: templates,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	err := h.AuthService.Register(user)
	if err != nil {
		http.Error(w, "Cannot register user", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// authenticate a user
	loggedUser, err := h.AuthService.Login(user.Email, user.Password)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	// Create session
	sessionID, err := h.SessionService.StartSession(loggedUser.ID)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}
	// 3. Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // enable in HTTPS
		Expires: time.Now().Add(24 * time.Hour),
	})
	json.NewEncoder(w).Encode(loggedUser)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.SessionService.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.Write([]byte("logged out"))
}