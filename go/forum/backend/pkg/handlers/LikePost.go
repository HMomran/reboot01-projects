package handlers

import (
	"forum/backend/database"
	"net/http"
)

func LikePost(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	handleReaction(w, r, true, DB)
}
