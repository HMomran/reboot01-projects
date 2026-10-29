package handlers

import (
	"fmt"
	"forum/backend"
	"forum/backend/database"
	"net/http"
)

func handleReaction(w http.ResponseWriter, r *http.Request, isLike bool, DB *database.DB) {
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

	postID := r.FormValue("post_id")
	if postID == "" {
		backend.RenderErrorPage(w, http.StatusBadRequest, "Post ID required")
		return
	}

	var postIDInt int
	_, err = fmt.Sscanf(postID, "%d", &postIDInt)
	if err != nil {
		backend.RenderErrorPage(w, http.StatusBadRequest, "Invalid post ID")
		return
	}

	err = DB.AddReaction(postIDInt, userID, isLike)
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error adding reaction")
		return
	}

	http.Redirect(w, r, "/communication", http.StatusSeeOther)
}
