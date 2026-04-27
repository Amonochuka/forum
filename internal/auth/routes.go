package auth

import (
	"fmt"
	"net/http"
)

func RegisterRoutes(handler *Handler) {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "web/templates/register.html")
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.Register(w, r)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "web/templates/login.html")
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.Login(w, r)
	})

	http.HandleFunc("/logout", handler.Logout)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// 1. Create a data container for the template
		data := make(map[string]interface{})

		// 2. Check for the session cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			data["User"] = nil
			fmt.Println("No cookie found")
		} else {
			// Use the function you already wrote!
			userID, err := handler.SessionService.ValidateSession(cookie.Value)
			if err != nil {
				http.SetCookie(w, &http.Cookie{Name: "session_id", MaxAge: -1})
				data["User"] = nil
				fmt.Println("Session invalid:", err)
			} else {
				// Now get the user details from your Auth/User service
				user, err := handler.AuthService.GetUserByID(userID)
				if err != nil {
					fmt.Println("User not found in DB:", err)
				} else {
					fmt.Printf("User found: %+v\n", user)
					data["User"] = user
				}
			}
		}

		handler.templates.ExecuteTemplate(w, "index.html", nil)
	})
}
