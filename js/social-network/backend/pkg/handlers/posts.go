package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"social-network/pkg/auth"

	"github.com/google/uuid"
)

type CreatePostRequest struct {
	Content        string   `json:"content"`
	Image          string   `json:"image"`
	Privacy        string   `json:"privacy"`
	AllowedUserIDs []string `json:"allowed_user_ids"`
}

type CreateCommentRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

type PostResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	Privacy   string `json:"privacy"`
	CreatedAt string `json:"created_at"`
}

type CommentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	Likes     int    `json:"likes"`
	ILiked    bool   `json:"i_liked"`
}

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
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

		var body CreatePostRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.Content == "" && body.Image == "" {
			http.Error(w, "Post must have text or an image", http.StatusBadRequest)
			return
		}
		if len(strings.TrimSpace(body.Content)) > 2000 {
			http.Error(w, "Post content cannot exceed 2000 characters", http.StatusBadRequest)
			return
		}

		if body.Privacy != "public" && body.Privacy != "followers" && body.Privacy != "private" {
			http.Error(w, "privacy must be 'public', 'followers', or 'private'", http.StatusBadRequest)
			return
		}

		postID := uuid.New().String()

		var imageValue interface{}
		if body.Image == "" {
			imageValue = nil
		} else {
			imageValue = body.Image
		}

		_, err = db.Exec(`
			INSERT INTO posts (id, user_id, content, image, privacy)
			VALUES (?, ?, ?, ?, ?)
		`, postID, currentUserID, body.Content, imageValue, body.Privacy)
		if err != nil {
			http.Error(w, "Could not create post", http.StatusInternalServerError)
			return
		}

		if body.Privacy == "private" {
			for _, allowedUserID := range body.AllowedUserIDs {
				if allowedUserID == "" {
					continue
				}

				var count int
				err = db.QueryRow(`
					SELECT COUNT(*)
					FROM followers
					WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
				`, allowedUserID, currentUserID).Scan(&count)
				if err != nil || count == 0 {
					continue
				}

				_, err = db.Exec(`
					INSERT INTO post_allowed_users (id, post_id, user_id)
					VALUES (?, ?, ?)
				`, uuid.New().String(), postID, allowedUserID)
				if err != nil {
					http.Error(w, "Could not save allowed users", http.StatusInternalServerError)
					return
				}
			}
		}

		w.Write([]byte("Post created successfully"))
	}
}

func ListVisiblePostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, _ := auth.GetCurrentUserID(db, r)

		rows, err := db.Query(`
			SELECT p.id, p.user_id, p.content, p.image, p.privacy, p.created_at
			FROM posts p
			WHERE
				p.privacy = 'public'
				OR p.user_id = ?
				OR (
					p.privacy = 'followers'
					AND EXISTS (
						SELECT 1
						FROM followers f
						WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
					)
				)
				OR (
					p.privacy = 'private'
					AND EXISTS (
						SELECT 1
						FROM post_allowed_users pau
						WHERE pau.post_id = p.id AND pau.user_id = ?
					)
				)
			ORDER BY p.created_at DESC
		`, currentUserID, currentUserID, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		posts := make([]PostResponse, 0)

		for rows.Next() {
			var item PostResponse
			var image sql.NullString

			err := rows.Scan(
				&item.ID,
				&item.UserID,
				&item.Content,
				&image,
				&item.Privacy,
				&item.CreatedAt,
			)
			if err != nil {
				http.Error(w, "Could not read posts", http.StatusInternalServerError)
				return
			}

			if image.Valid {
				item.Image = image.String
			} else {
				item.Image = ""
			}

			posts = append(posts, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func CreateCommentHandler(db *sql.DB) http.HandlerFunc {
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

		var body CreateCommentRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.PostID == "" || body.Content == "" {
			http.Error(w, "post_id and content are required", http.StatusBadRequest)
			return
		}
		if len(strings.TrimSpace(body.Content)) > 500 {
			http.Error(w, "Comment cannot exceed 500 characters", http.StatusBadRequest)
			return
		}

		canView, err := canUserViewPost(db, currentUserID, body.PostID)
		if err != nil {
			http.Error(w, "Could not verify post access", http.StatusInternalServerError)
			return
		}
		if !canView {
			http.Error(w, "You cannot comment on this post", http.StatusForbidden)
			return
		}

		commentID := uuid.New().String()

		var imageValue interface{}
		if body.Image == "" {
			imageValue = nil
		} else {
			imageValue = body.Image
		}

		_, err = db.Exec(`
			INSERT INTO comments (id, post_id, user_id, content, image)
			VALUES (?, ?, ?, ?, ?)
		`, commentID, body.PostID, currentUserID, body.Content, imageValue)
		if err != nil {
			http.Error(w, "Could not create comment", http.StatusInternalServerError)
			return
		}

		// Notify the post author and broadcast the new comment count to feed viewers
		var postOwner, firstName, lastName string
		db.QueryRow("SELECT user_id FROM posts WHERE id = ?", body.PostID).Scan(&postOwner)
		db.QueryRow("SELECT first_name, last_name FROM users WHERE id = ?", currentUserID).Scan(&firstName, &lastName)
		if postOwner != "" && postOwner != currentUserID {
			_ = createNotification(db, postOwner, "comment", body.PostID,
				firstName+" "+lastName+" commented on your post")
		}
		if GlobalFeedHub != nil {
			var commentCount int
			db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", body.PostID).Scan(&commentCount)
			GlobalFeedHub.Broadcast(map[string]interface{}{
				"event":         "post_updated",
				"post_id":       body.PostID,
				"comment_count": commentCount,
			})
		}

		w.Write([]byte("Comment created successfully"))
	}
}

func ListPostCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		postID := r.URL.Query().Get("post_id")
		if postID == "" {
			http.Error(w, "Missing post_id", http.StatusBadRequest)
			return
		}

		currentUserID, _ := auth.GetCurrentUserID(db, r)

		canView, err := canUserViewPost(db, currentUserID, postID)
		if err != nil {
			http.Error(w, "Could not verify post access", http.StatusInternalServerError)
			return
		}
		if !canView {
			http.Error(w, "You cannot view comments for this post", http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT c.id, c.post_id, c.user_id,
			       u.first_name, u.last_name, COALESCE(u.avatar, ''),
			       c.content, COALESCE(c.image, ''), c.created_at,
			       COUNT(cl.id),
			       COALESCE(SUM(CASE WHEN cl.user_id = ? THEN 1 ELSE 0 END), 0)
			FROM comments c
			JOIN users u ON u.id = c.user_id
			LEFT JOIN comment_likes cl ON cl.comment_id = c.id
			WHERE c.post_id = ?
			GROUP BY c.id, c.post_id, c.user_id, u.first_name, u.last_name, u.avatar, c.content, c.image, c.created_at
			ORDER BY c.created_at ASC
		`, currentUserID, postID)
		if err != nil {
			http.Error(w, "Could not fetch comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		comments := make([]CommentResponse, 0)

		for rows.Next() {
			var item CommentResponse
			var iLiked int

			err := rows.Scan(
				&item.ID, &item.PostID, &item.UserID,
				&item.FirstName, &item.LastName, &item.Avatar,
				&item.Content, &item.Image, &item.CreatedAt,
				&item.Likes, &iLiked,
			)
			if err != nil {
				http.Error(w, "Could not read comments", http.StatusInternalServerError)
				return
			}
			item.ILiked = iLiked > 0

			comments = append(comments, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

func canUserViewPost(db *sql.DB, currentUserID string, postID string) (bool, error) {
	var ownerID string
	var privacy string

	err := db.QueryRow(`
		SELECT user_id, privacy
		FROM posts
		WHERE id = ?
	`, postID).Scan(&ownerID, &privacy)
	if err != nil {
		return false, err
	}

	if privacy == "public" {
		return true, nil
	}

	if currentUserID == ownerID {
		return true, nil
	}

	if privacy == "followers" {
		var count int
		err = db.QueryRow(`
			SELECT COUNT(*)
			FROM followers
			WHERE follower_id = ? AND following_id = ? AND status = 'accepted'
		`, currentUserID, ownerID).Scan(&count)
		if err != nil {
			return false, err
		}
		return count > 0, nil
	}

	if privacy == "private" {
		var count int
		err = db.QueryRow(`
			SELECT COUNT(*)
			FROM post_allowed_users
			WHERE post_id = ? AND user_id = ?
		`, postID, currentUserID).Scan(&count)
		if err != nil {
			return false, err
		}
		return count > 0, nil
	}

	return false, nil
}

func ExploreHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		rows, err := db.Query(`
			SELECT p.id, p.user_id, u.first_name, u.last_name, COALESCE(u.avatar, ''),
			       p.content, COALESCE(p.image, ''), p.created_at,
			       COALESCE(SUM(CASE WHEN pr.type = 'like' THEN 1 ELSE 0 END), 0),
			       COALESCE(SUM(CASE WHEN pr.type = 'dislike' THEN 1 ELSE 0 END), 0),
			       COALESCE(MAX(CASE WHEN pr.user_id = ? THEN pr.type ELSE NULL END), '')
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN post_reactions pr ON pr.post_id = p.id
			WHERE p.privacy = 'public'
			  AND u.profile_type = 'public'
			GROUP BY p.id, p.user_id, u.first_name, u.last_name, u.avatar, p.content, p.image, p.created_at
			ORDER BY p.created_at DESC
		`, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch explore posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		posts := make([]FeedPostResponse, 0)
		for rows.Next() {
			var item FeedPostResponse
			err := rows.Scan(
				&item.ID, &item.UserID, &item.FirstName, &item.LastName, &item.Avatar,
				&item.Content, &item.Image, &item.CreatedAt,
				&item.Likes, &item.Dislikes, &item.MyReaction,
			)
			if err != nil {
				http.Error(w, "Could not read explore posts", http.StatusInternalServerError)
				return
			}
			item.Privacy = "public"
			posts = append(posts, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

type FeedPostResponse struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Avatar       string `json:"avatar"`
	Content      string `json:"content"`
	Image        string `json:"image"`
	Privacy      string `json:"privacy"`
	CreatedAt    string `json:"created_at"`
	Likes        int    `json:"likes"`
	Dislikes     int    `json:"dislikes"`
	MyReaction   string `json:"my_reaction"`
	CommentCount int    `json:"comment_count"`
}

func FeedHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, err := auth.GetCurrentUserID(db, r)
		var rows *sql.Rows
		var errQuery error

		if err != nil {
			rows, errQuery = db.Query(`
				SELECT p.id, p.user_id, u.first_name, u.last_name, COALESCE(u.avatar, ''),
				       p.content, COALESCE(p.image, ''), p.privacy, p.created_at,
				       COALESCE(SUM(CASE WHEN pr.type = 'like' THEN 1 ELSE 0 END), 0),
				       COALESCE(SUM(CASE WHEN pr.type = 'dislike' THEN 1 ELSE 0 END), 0),
				       '',
				       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id)
				FROM posts p
				JOIN users u ON u.id = p.user_id
				LEFT JOIN post_reactions pr ON pr.post_id = p.id
				WHERE p.privacy = 'public'
				GROUP BY p.id, p.user_id, u.first_name, u.last_name, u.avatar, p.content, p.image, p.privacy, p.created_at
				ORDER BY p.created_at DESC
			`)
		} else {
			rows, errQuery = db.Query(`
				SELECT p.id, p.user_id, u.first_name, u.last_name, COALESCE(u.avatar, ''),
				       p.content, COALESCE(p.image, ''), p.privacy, p.created_at,
				       COALESCE(SUM(CASE WHEN pr.type = 'like' THEN 1 ELSE 0 END), 0),
				       COALESCE(SUM(CASE WHEN pr.type = 'dislike' THEN 1 ELSE 0 END), 0),
				       COALESCE(MAX(CASE WHEN pr.user_id = ? THEN pr.type ELSE NULL END), ''),
				       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id)
				FROM posts p
				JOIN users u ON u.id = p.user_id
				LEFT JOIN post_reactions pr ON pr.post_id = p.id
				WHERE
					-- Your own posts are always visible
					p.user_id = ?
					-- Public posts are visible to everyone who is logged in
					OR p.privacy = 'public'
					-- Followers-only posts: visible if you follow the author (accepted)
					OR (
						p.privacy = 'followers'
						AND EXISTS (
							SELECT 1 FROM followers f
							WHERE f.follower_id = ? AND f.following_id = p.user_id AND f.status = 'accepted'
						)
					)
					-- Private posts: visible only if you are on the allowlist
					OR (
						p.privacy = 'private'
						AND EXISTS (
							SELECT 1 FROM post_allowed_users pau
							WHERE pau.post_id = p.id AND pau.user_id = ?
						)
					)
				GROUP BY p.id, p.user_id, u.first_name, u.last_name, u.avatar, p.content, p.image, p.privacy, p.created_at
				ORDER BY p.created_at DESC
			`, currentUserID, currentUserID, currentUserID, currentUserID)
		}

		if errQuery != nil {
			http.Error(w, "Could not fetch feed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		posts := make([]FeedPostResponse, 0)
		for rows.Next() {
			var item FeedPostResponse
			err := rows.Scan(
				&item.ID, &item.UserID, &item.FirstName, &item.LastName, &item.Avatar,
				&item.Content, &item.Image, &item.Privacy, &item.CreatedAt,
				&item.Likes, &item.Dislikes, &item.MyReaction, &item.CommentCount,
			)
			if err != nil {
				http.Error(w, "Could not read feed", http.StatusInternalServerError)
				return
			}
			posts = append(posts, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

// UserPostsHandler returns all posts for a given user, respecting privacy.
// GET /profile/posts?user_id=X
func UserPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		targetUserID := r.URL.Query().Get("user_id")
		if targetUserID == "" {
			http.Error(w, "Missing user_id", http.StatusBadRequest)
			return
		}

		var profileType string
		err = db.QueryRow("SELECT profile_type FROM users WHERE id = ?", targetUserID).Scan(&profileType)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// If not the owner, check if the requester may view posts
		if currentUserID != targetUserID && profileType == "private" {
			var count int
			_ = db.QueryRow(
				"SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ? AND status = 'accepted'",
				currentUserID, targetUserID,
			).Scan(&count)
			if count == 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"error": "private"})
				return
			}
		}

		rows, err := db.Query(`
			SELECT p.id, p.user_id, u.first_name, u.last_name, COALESCE(u.avatar, ''),
			       p.content, COALESCE(p.image, ''), p.privacy, p.created_at,
			       COALESCE(SUM(CASE WHEN pr.type = 'like' THEN 1 ELSE 0 END), 0),
			       COALESCE(SUM(CASE WHEN pr.type = 'dislike' THEN 1 ELSE 0 END), 0),
			       COALESCE(MAX(CASE WHEN pr.user_id = ? THEN pr.type ELSE NULL END), ''),
			       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id)
			FROM posts p
			JOIN users u ON u.id = p.user_id
			LEFT JOIN post_reactions pr ON pr.post_id = p.id
			WHERE p.user_id = ?
			GROUP BY p.id, p.user_id, u.first_name, u.last_name, u.avatar, p.content, p.image, p.privacy, p.created_at
			ORDER BY p.created_at DESC
		`, currentUserID, targetUserID)
		if err != nil {
			http.Error(w, "Could not fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		userPosts := make([]FeedPostResponse, 0)
		for rows.Next() {
			var item FeedPostResponse
			err := rows.Scan(
				&item.ID, &item.UserID, &item.FirstName, &item.LastName, &item.Avatar,
				&item.Content, &item.Image, &item.Privacy, &item.CreatedAt,
				&item.Likes, &item.Dislikes, &item.MyReaction, &item.CommentCount,
			)
			if err != nil {
				continue
			}
			userPosts = append(userPosts, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userPosts)
	}
}

func DeletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		postID := r.URL.Query().Get("post_id")
		if postID == "" {
			http.Error(w, "post_id is required", http.StatusBadRequest)
			return
		}
		var ownerID string
		err = db.QueryRow(`SELECT user_id FROM posts WHERE id = ?`, postID).Scan(&ownerID)
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if ownerID != currentUserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		_, _ = tx.Exec(`DELETE FROM comment_likes WHERE comment_id IN (SELECT id FROM comments WHERE post_id = ?)`, postID)
		_, _ = tx.Exec(`DELETE FROM comments WHERE post_id = ?`, postID)
		_, _ = tx.Exec(`DELETE FROM post_reactions WHERE post_id = ?`, postID)
		_, _ = tx.Exec(`DELETE FROM post_allowed_users WHERE post_id = ?`, postID)
		_, err = tx.Exec(`DELETE FROM posts WHERE id = ? AND user_id = ?`, postID, currentUserID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Could not delete post", http.StatusInternalServerError)
			return
		}
		tx.Commit()
		w.Write([]byte("Post deleted"))
	}
}

func DeleteCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		commentID := r.URL.Query().Get("comment_id")
		if commentID == "" {
			http.Error(w, "comment_id is required", http.StatusBadRequest)
			return
		}
		var ownerID, postID string
		err = db.QueryRow(`SELECT user_id, post_id FROM comments WHERE id = ?`, commentID).Scan(&ownerID, &postID)
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		// Allow comment author OR post author to delete
		var postAuthorID string
		_ = db.QueryRow(`SELECT user_id FROM posts WHERE id = ?`, postID).Scan(&postAuthorID)
		if ownerID != currentUserID && postAuthorID != currentUserID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		_, _ = db.Exec(`DELETE FROM comment_likes WHERE comment_id = ?`, commentID)
		_, err = db.Exec(`DELETE FROM comments WHERE id = ?`, commentID)
		if err != nil {
			http.Error(w, "Could not delete comment", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Comment deleted"))
	}
}
