package pages

import (
	"forum/backend/database"
	"forum/backend"
	"html/template"
	"net/http"
)

func WelcomePage(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	if r.URL.Path != "/" {
		backend.RenderErrorPage(w, http.StatusNotFound, "The page you're looking for doesn't exist.")
		return
	}

	cookie, err := r.Cookie("session_token")
	if err == nil {
		_, err := DB.GetUserIDFromSession(cookie.Value)
		if err == nil {
			http.Redirect(w, r, "/communication", http.StatusSeeOther)
			return
		}
	}

	tmpl, err := template.ParseFiles("Frontend/templates/WelcomePage.html")
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading welcome page.")
		return
	}
	tmpl.Execute(w, nil)
}