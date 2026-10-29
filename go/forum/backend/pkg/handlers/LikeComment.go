package handlers

import (
	"forum/backend/database"
	"net/http"
)

func LikeComment(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	handleCommentReaction(w, r, true, DB)
}