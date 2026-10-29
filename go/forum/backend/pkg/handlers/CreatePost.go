package handlers

import (
	"fmt"
	"forum/backend"
	"forum/backend/database"
	"net/http"
	"strings"
	"net/url"
)

func CreatePost(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	userID, err := DB.GetUserIDFromSession(cookie.Value)
	if err != nil {
		backend.RenderErrorPage(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/create-post-page", http.StatusSeeOther)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	r.ParseForm()
	categoryStrings := r.Form["categories"]

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	formData := FormData{
		Title:              title,
		Content:            content,
		SelectedCategories: categoryStrings,
	}

	if len(title) == 0 {
		errorMsg := url.QueryEscape("Post title is required.")
		http.Redirect(w, r, "/create-post-page?error="+errorMsg, http.StatusSeeOther)
		return
	}
	if len(title) > 20 {
		errorMsg := url.QueryEscape("Post title too long (max 20 characters).")
		http.Redirect(w, r, "/create-post-page?error="+errorMsg, http.StatusSeeOther)
		return
	}
	if len(content) == 0 {
		errorMsg := url.QueryEscape("Post content is required.")
		http.Redirect(w, r, "/create-post-page?error="+errorMsg, http.StatusSeeOther)
		return
	}
	if len(content) > 50 {
		errorMsg := url.QueryEscape("Post content too long (max 50 characters).")
		http.Redirect(w, r, "/create-post-page?error="+errorMsg, http.StatusSeeOther)
		return
	}
	if len(categoryStrings) == 0 {
		errorMsg := url.QueryEscape("Please select at least one category for your post.")
		http.Redirect(w, r, "/create-post-page?error="+errorMsg, http.StatusSeeOther)
		return
	}

	var categoryIDs []int
	for _, catStr := range categoryStrings {
		var catID int
		fmt.Sscanf(catStr, "%d", &catID)
		categoryIDs = append(categoryIDs, catID)
	}

	err = DB.CreatePostWithCategories(userID, title, content, categoryIDs)
	if err != nil {
		showCreatePostFormPage(w, userID, "Error creating post. Please try again.", formData, DB)
		return
	}

	http.Redirect(w, r, "/communication", http.StatusSeeOther)
}
