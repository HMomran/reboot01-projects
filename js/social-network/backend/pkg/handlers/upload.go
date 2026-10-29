package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"social-network/pkg/auth"
	"social-network/pkg/utils"
)

func UploadAvatarHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseMultipartForm(5 << 20)
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		path, err := utils.SaveUploadedImage(r, "avatar", "uploads/avatars")
		if err != nil {
			http.Error(w, "Could not save avatar: "+err.Error(), http.StatusBadRequest)
			return
		}

		publicPath := "/" + strings.ReplaceAll(path, "\\", "/")

		_, err = db.Exec(`
			UPDATE users
			SET avatar = ?
			WHERE id = ?
		`, publicPath, currentUserID)
		if err != nil {
			http.Error(w, "Could not update avatar", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(publicPath))
	}
}

func UploadPostImageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseMultipartForm(5 << 20)
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		path, err := utils.SaveUploadedImage(r, "image", "uploads/posts")
		if err != nil {
			http.Error(w, "Could not save image: "+err.Error(), http.StatusBadRequest)
			return
		}

		publicPath := "/" + strings.ReplaceAll(path, "\\", "/")
		w.Write([]byte(publicPath))
	}
}

func UploadCommentImageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseMultipartForm(5 << 20)
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		path, err := utils.SaveUploadedImage(r, "image", "uploads/comments")
		if err != nil {
			http.Error(w, "Could not save image: "+err.Error(), http.StatusBadRequest)
			return
		}

		publicPath := "/" + strings.ReplaceAll(path, "\\", "/")
		w.Write([]byte(publicPath))
	}
}

func UploadGroupPostImageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseMultipartForm(5 << 20)
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		path, err := utils.SaveUploadedImage(r, "image", "uploads/group_posts")
		if err != nil {
			http.Error(w, "Could not save image: "+err.Error(), http.StatusBadRequest)
			return
		}

		publicPath := "/" + strings.ReplaceAll(path, "\\", "/")
		w.Write([]byte(publicPath))
	}
}

func UploadGroupCommentImageHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		_, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseMultipartForm(5 << 20)
		if err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}

		path, err := utils.SaveUploadedImage(r, "image", "uploads/group_comments")
		if err != nil {
			http.Error(w, "Could not save image: "+err.Error(), http.StatusBadRequest)
			return
		}

		publicPath := "/" + strings.ReplaceAll(path, "\\", "/")
		w.Write([]byte(publicPath))
	}
}