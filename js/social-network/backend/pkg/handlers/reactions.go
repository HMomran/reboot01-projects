package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social-network/pkg/auth"

	"github.com/google/uuid"
)

type ReactToPostRequest struct {
	PostID string `json:"post_id"`
	Type   string `json:"type"` // "like", "dislike", or "" to remove
}

type ReactionResponse struct {
	Likes      int    `json:"likes"`
	Dislikes   int    `json:"dislikes"`
	MyReaction string `json:"my_reaction"`
}

// ReactToPostHandler handles POST /posts/react
// Body: { post_id, type: "like"|"dislike"|"" }
// Toggling the same type removes it; switching type changes it.
func ReactToPostHandler(db *sql.DB) http.HandlerFunc {
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

		var body ReactToPostRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.PostID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.Type != "" && body.Type != "like" && body.Type != "dislike" {
			http.Error(w, "type must be 'like', 'dislike', or empty", http.StatusBadRequest)
			return
		}

		var existing string
		err = db.QueryRow(
			"SELECT type FROM post_reactions WHERE post_id = ? AND user_id = ?",
			body.PostID, currentUserID,
		).Scan(&existing)

		switch err {
		case sql.ErrNoRows:
			// No existing reaction
			if body.Type != "" {
				_, err = db.Exec(
					"INSERT INTO post_reactions (id, post_id, user_id, type) VALUES (?, ?, ?, ?)",
					uuid.New().String(), body.PostID, currentUserID, body.Type,
				)
				if err != nil {
					http.Error(w, "Could not save reaction", http.StatusInternalServerError)
					return
				}
			}
		case nil:
			if body.Type == "" || existing == body.Type {
				// Toggle off
				_, err = db.Exec(
					"DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?",
					body.PostID, currentUserID,
				)
			} else {
				// Switch reaction type
				_, err = db.Exec(
					"UPDATE post_reactions SET type = ? WHERE post_id = ? AND user_id = ?",
					body.Type, body.PostID, currentUserID,
				)
			}
			if err != nil {
				http.Error(w, "Could not update reaction", http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Could not read reaction", http.StatusInternalServerError)
			return
		}

		// Return updated counts
		resp := ReactionResponse{}
		db.QueryRow(
			"SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND type = 'like'",
			body.PostID,
		).Scan(&resp.Likes)
		db.QueryRow(
			"SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND type = 'dislike'",
			body.PostID,
		).Scan(&resp.Dislikes)
		db.QueryRow(
			"SELECT COALESCE(type, '') FROM post_reactions WHERE post_id = ? AND user_id = ?",
			body.PostID, currentUserID,
		).Scan(&resp.MyReaction)

		// Notify the post author (only on a new like, not removal/dislike)
		if body.Type == "like" {
			var postOwner, firstName, lastName string
			db.QueryRow("SELECT user_id FROM posts WHERE id = ?", body.PostID).Scan(&postOwner)
			db.QueryRow("SELECT first_name, last_name FROM users WHERE id = ?", currentUserID).Scan(&firstName, &lastName)
			if postOwner != "" && postOwner != currentUserID {
				_ = createNotification(db, postOwner, "like", body.PostID,
					firstName+" "+lastName+" liked your post")
			}
		}

		// Broadcast feed update to all live feed viewers
		if GlobalFeedHub != nil {
			GlobalFeedHub.Broadcast(map[string]interface{}{
				"event":    "post_updated",
				"post_id":  body.PostID,
				"likes":    resp.Likes,
				"dislikes": resp.Dislikes,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

type LikeCommentRequest struct {
	CommentID string `json:"comment_id"`
}

type CommentLikeResponse struct {
	Likes  int  `json:"likes"`
	ILiked bool `json:"i_liked"`
}

// LikeCommentHandler handles POST /posts/comments/like
// Toggles a like on a comment.
func LikeCommentHandler(db *sql.DB) http.HandlerFunc {
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

		var body LikeCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.CommentID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var existingID string
		err = db.QueryRow(
			"SELECT id FROM comment_likes WHERE comment_id = ? AND user_id = ?",
			body.CommentID, currentUserID,
		).Scan(&existingID)

		switch err {
		case sql.ErrNoRows:
			_, err = db.Exec(
				"INSERT INTO comment_likes (id, comment_id, user_id) VALUES (?, ?, ?)",
				uuid.New().String(), body.CommentID, currentUserID,
			)
		case nil:
			_, err = db.Exec(
				"DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?",
				body.CommentID, currentUserID,
			)
		default:
			http.Error(w, "Could not read comment like", http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, "Could not update comment like", http.StatusInternalServerError)
			return
		}

		resp := CommentLikeResponse{}
		db.QueryRow(
			"SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?",
			body.CommentID,
		).Scan(&resp.Likes)
		var cnt int
		db.QueryRow(
			"SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND user_id = ?",
			body.CommentID, currentUserID,
		).Scan(&cnt)
		resp.ILiked = cnt > 0

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
