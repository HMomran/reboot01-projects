package pages

import (
	"forum/backend"
	"forum/backend/database"
	"html/template"
	"net/http"
)

type FormData struct {
	Title              string
	Content            string
	SelectedCategories []string
}

func CreatePostPage(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	userID, err := DB.GetUserIDFromSession(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	categories, err := DB.GetAllCategories()
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading categories")
		return
	}
	errorMsg := r.URL.Query().Get("error")
	data := struct {
		Categories    []database.Category
		UserID        int
		CategoryError string
		FormData      FormData
		ErrorMsg         string
	}{
		Categories:    categories,
		UserID:        userID,
		CategoryError: "",
		FormData:      FormData{},
		ErrorMsg:         errorMsg,
	}

	tmpl, err := template.ParseFiles("Frontend/templates/CreatePost.html")
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading page")
		return
	}
	tmpl.Execute(w, data)
}