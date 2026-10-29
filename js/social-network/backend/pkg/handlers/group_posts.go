package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"social-network/pkg/auth"

	"github.com/google/uuid"
)

type CreateGroupPostRequest struct {
	GroupID string `json:"group_id"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

type GroupPostResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
}

type CreateGroupCommentRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

type GroupCommentResponse struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Content   string `json:"content"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
}

func CreateGroupPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var body CreateGroupPostRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.GroupID == "" || body.Content == "" {
			http.Error(w, "group_id and content required", http.StatusBadRequest)
			return
		}
		if len(strings.TrimSpace(body.Content)) > 2000 {
			http.Error(w, "Post content cannot exceed 2000 characters", http.StatusBadRequest)
			return
		}

		if !isGroupMember(db, body.GroupID, currentUserID) {
			http.Error(w, "Not a group member", http.StatusForbidden)
			return
		}

		postID := uuid.New().String()

		_, err = db.Exec(`
			INSERT INTO group_posts (id, group_id, user_id, content, image)
			VALUES (?, ?, ?, ?, ?)
		`, postID, body.GroupID, currentUserID, body.Content, nullIfEmptyGroup(body.Image))

		if err != nil {
			http.Error(w, "Could not create group post", 500)
			return
		}

		w.Write([]byte("Group post created"))
	}
}

func ListGroupPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "group_id required", http.StatusBadRequest)
			return
		}

		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !isGroupMember(db, groupID, currentUserID) {
			http.Error(w, "Not a group member", http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT gp.id, gp.group_id, gp.user_id,
			       u.first_name, u.last_name, COALESCE(u.avatar, ''),
			       gp.content, COALESCE(gp.image, ''), gp.created_at
			FROM group_posts gp
			JOIN users u ON u.id = gp.user_id
			WHERE gp.group_id = ?
			ORDER BY gp.created_at DESC
			LIMIT 50
		`, groupID)
		if err != nil {
			http.Error(w, "Could not fetch group posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		posts := make([]GroupPostResponse, 0)

		for rows.Next() {
			var p GroupPostResponse
			err := rows.Scan(&p.ID, &p.GroupID, &p.UserID, &p.FirstName, &p.LastName, &p.Avatar, &p.Content, &p.Image, &p.CreatedAt)
			if err != nil {
				http.Error(w, "Could not read group posts", http.StatusInternalServerError)
				return
			}
			posts = append(posts, p)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

func CreateGroupCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var body CreateGroupCommentRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.PostID == "" || body.Content == "" {
			http.Error(w, "post_id and content required", 400)
			return
		}
		if len(strings.TrimSpace(body.Content)) > 500 {
			http.Error(w, "Comment cannot exceed 500 characters", http.StatusBadRequest)
			return
		}

		groupID := getGroupIDFromPost(db, body.PostID)

		if !isGroupMember(db, groupID, currentUserID) {
			http.Error(w, "Not a group member", 403)
			return
		}

		commentID := uuid.New().String()

		_, err = db.Exec(`
			INSERT INTO group_comments (id, post_id, user_id, content, image)
			VALUES (?, ?, ?, ?, ?)
		`, commentID, body.PostID, currentUserID, body.Content, nullIfEmptyGroup(body.Image))

		if err != nil {
			http.Error(w, "Could not create comment", 500)
			return
		}

		w.Write([]byte("Group comment created"))
	}
}

func ListGroupCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		postID := r.URL.Query().Get("post_id")
		if postID == "" {
			http.Error(w, "post_id required", http.StatusBadRequest)
			return
		}

		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		groupID := getGroupIDFromPost(db, postID)
		if groupID == "" {
			http.Error(w, "Group post not found", http.StatusNotFound)
			return
		}

		if !isGroupMember(db, groupID, currentUserID) {
			http.Error(w, "Not allowed", http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT gc.id, gc.post_id, gc.user_id,
			       u.first_name, u.last_name, COALESCE(u.avatar, ''),
			       gc.content, COALESCE(gc.image, ''), gc.created_at
			FROM group_comments gc
			JOIN users u ON u.id = gc.user_id
			WHERE gc.post_id = ?
			ORDER BY gc.created_at ASC
		`, postID)
		if err != nil {
			http.Error(w, "Could not fetch group comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		comments := make([]GroupCommentResponse, 0)

		for rows.Next() {
			var c GroupCommentResponse
			err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.FirstName, &c.LastName, &c.Avatar, &c.Content, &c.Image, &c.CreatedAt)
			if err != nil {
				http.Error(w, "Could not read group comments", http.StatusInternalServerError)
				return
			}
			comments = append(comments, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}

func DeleteGroupPostHandler(db *sql.DB) http.HandlerFunc {
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
			http.Error(w, "post_id required", http.StatusBadRequest)
			return
		}

		var postAuthorID, groupID string
		err = db.QueryRow(`SELECT user_id, group_id FROM group_posts WHERE id = ?`, postID).Scan(&postAuthorID, &groupID)
		if err != nil {
			http.Error(w, "Group post not found", http.StatusNotFound)
			return
		}

		var groupCreatorID string
		_ = db.QueryRow(`SELECT creator_id FROM groups WHERE id = ?`, groupID).Scan(&groupCreatorID)

		if currentUserID != postAuthorID && currentUserID != groupCreatorID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		_, _ = tx.Exec(`DELETE FROM group_comments WHERE post_id = ?`, postID)
		_, err = tx.Exec(`DELETE FROM group_posts WHERE id = ?`, postID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Could not delete group post", http.StatusInternalServerError)
			return
		}
		tx.Commit()
		w.Write([]byte("Group post deleted"))
	}
}

func DeleteGroupCommentHandler(db *sql.DB) http.HandlerFunc {
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
			http.Error(w, "comment_id required", http.StatusBadRequest)
			return
		}

		var commentAuthorID string
		err = db.QueryRow(`SELECT user_id FROM group_comments WHERE id = ?`, commentID).Scan(&commentAuthorID)
		if err != nil {
			http.Error(w, "Group comment not found", http.StatusNotFound)
			return
		}
		if currentUserID != commentAuthorID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		_, err = db.Exec(`DELETE FROM group_comments WHERE id = ?`, commentID)
		if err != nil {
			http.Error(w, "Could not delete group comment", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Group comment deleted"))
	}
}

func isGroupMember(db *sql.DB, groupID, userID string) bool {
	var count int
	db.QueryRow(`
		SELECT COUNT(*)
		FROM group_members
		WHERE group_id = ? AND user_id = ? AND status = 'accepted'
	`, groupID, userID).Scan(&count)
	return count > 0
}

func getGroupIDFromPost(db *sql.DB, postID string) string {
	var groupID string
	db.QueryRow(`
		SELECT group_id FROM group_posts WHERE id = ?
	`, postID).Scan(&groupID)
	return groupID
}

func nullIfEmptyGroup(val string) interface{} {
	if val == "" {
		return nil
	}
	return val
}
