package handlers

import (
	"fmt"
	"forum/backend"
	"forum/backend/database"
	"net/http"
)


func handleCommentReaction(w http.ResponseWriter, r *http.Request, isLike bool, DB *database.DB) {
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

	commentID := r.FormValue("comment_id")
	if commentID == "" {
		backend.RenderErrorPage(w, http.StatusBadRequest, "Comment ID required")
		return
	}

	var commentIDInt int
	_, err = fmt.Sscanf(commentID, "%d", &commentIDInt)
	if err != nil {
		backend.RenderErrorPage(w, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	err = DB.AddCommentReaction(commentIDInt, userID, isLike)
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error adding reaction")
		return
	}

	http.Redirect(w, r, "/communication", http.StatusSeeOther)
}
