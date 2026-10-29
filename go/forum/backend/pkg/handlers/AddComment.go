package handlers

import (
	"fmt"
	"forum/backend"
	"forum/backend/database"
	"net/http"
	"net/url"
	"strings"
)

func AddComment(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	if r.Method != http.MethodPost {
		backend.RenderErrorPage(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		backend.RenderErrorPage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := DB.GetUserIDFromSession(cookie.Value)
	if err != nil {
		backend.RenderErrorPage(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")
	content = strings.TrimSpace(content)

	if content == "" {
		errorMsg := url.QueryEscape("Comment cannot be empty or contain only spaces")
		http.Redirect(w, r, "/communication?error="+errorMsg, http.StatusSeeOther)
		return
	}
	if len(content) > 50 {
		errorMsg := url.QueryEscape("Comment is too long (max 50 characters)")
		http.Redirect(w, r, "/communication?error="+errorMsg, http.StatusSeeOther)
		return
	}

	var postID int
	_, err = fmt.Sscanf(postIDStr, "%d", &postID)
	if err != nil {
		backend.RenderErrorPage(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = DB.CreateComment(postID, userID, content)
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error adding comment")
		return
	}

	http.Redirect(w, r, "/communication", http.StatusSeeOther)
}
