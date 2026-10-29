package handlers

import (
	"forum/backend/database"
	"net/http"
)

func DislikeComment(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	handleCommentReaction(w, r, false, DB)
}