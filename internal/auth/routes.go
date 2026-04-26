package auth

import (
	"html/template"
	"net/http"
)

func RegisterRoutes(handler *Handler) {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "web/templates/register.html")
			return
		}
		handler.Register(w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "web/templates/login.html")
			return
		}
		handler.Login(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmpl, err := template.ParseFiles(
			"web/templates/index.html",
			"web/templates/components/navbar.html",
			"web/templates/components/hero.html",
			"web/templates/components/create_post.html",
			"web/templates/components/sidebar.html",
			"web/templates/components/footer.html",
			"web/templates/components/scripts.html",
		)
		if err != nil {
			http.Error(w, "Failed to load templates: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})
}