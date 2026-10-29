package handlers

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

func showCreatePostFormPage(w http.ResponseWriter, userID int, errorMsg string, formData FormData, DB *database.DB) {
	categories, err := DB.GetAllCategories()
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading categories")
		return
	}

	data := struct {
		Categories    []database.Category
		UserID        int
		CategoryError string
		FormData      FormData
	}{
		Categories:    categories,
		UserID:        userID,
		CategoryError: errorMsg,
		FormData:      formData,
	}

	tmpl, err := template.ParseFiles("Frontend/templates/CreatePost.html")
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading page")
		return
	}
	tmpl.Execute(w, data)
}
