package handlers

import (
	"forum/backend/database"
	"net/http"
)

func DislikePost(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	handleReaction(w, r, false, DB)
}
