package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social-network/pkg/auth"
	"social-network/pkg/utils"
)

type PrivateMessageResponse struct {
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	IsRead     bool   `json:"is_read"`
}

type GroupMessageResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// UnreadConversation summarises unread messages from one sender.
type UnreadConversation struct {
	SenderID    string `json:"sender_id"`
	Count       int    `json:"count"`
	LastMessage string `json:"last_message"`
	LastAt      string `json:"last_at"`
}

func ListPrivateMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		otherUserID := r.URL.Query().Get("user_id")
		if otherUserID == "" {
			http.Error(w, "user_id required", http.StatusBadRequest)
			return
		}

		if !canUsersChatHandlers(db, currentUserID, otherUserID) {
			http.Error(w, "Not allowed to chat with this user", http.StatusForbidden)
			return
		}

		// Optional cursor: load messages sent before this timestamp (for pagination)
		before := r.URL.Query().Get("before")

		var rows *sql.Rows
		if before != "" {
			rows, err = db.Query(`
				SELECT id, sender_id, receiver_id, content, created_at, is_read
				FROM private_messages
				WHERE ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
				  AND created_at < ?
				ORDER BY created_at DESC
				LIMIT 10
			`, currentUserID, otherUserID, otherUserID, currentUserID, before)
		} else {
			rows, err = db.Query(`
				SELECT id, sender_id, receiver_id, content, created_at, is_read
				FROM (
					SELECT id, sender_id, receiver_id, content, created_at, is_read
					FROM private_messages
					WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
					ORDER BY created_at DESC
					LIMIT 10
				) sub
				ORDER BY created_at ASC
			`, currentUserID, otherUserID, otherUserID, currentUserID)
		}
		if err != nil {
			http.Error(w, "Could not fetch messages", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		messages := make([]PrivateMessageResponse, 0)
		for rows.Next() {
			var m PrivateMessageResponse
			var isRead int
			if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt, &isRead); err != nil {
				http.Error(w, "Could not read messages", http.StatusInternalServerError)
				return
			}
			m.IsRead = isRead == 1
			messages = append(messages, m)
		}

		// When paginating backwards, the DB returned DESC; reverse to chronological order
		if before != "" {
			for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
				messages[i], messages[j] = messages[j], messages[i]
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

// GetUnreadCountHandler returns unread message counts grouped by sender conversation.
// GET /messages/unread-count
func GetUnreadCountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		rows, err := db.Query(`
			SELECT
				sender_id,
				COUNT(*) AS cnt,
				(SELECT content FROM private_messages pm2
				 WHERE pm2.sender_id = pm.sender_id AND pm2.receiver_id = ?
				 ORDER BY pm2.created_at DESC LIMIT 1) AS last_msg,
				MAX(created_at) AS last_at
			FROM private_messages pm
			WHERE receiver_id = ? AND is_read = 0
			GROUP BY sender_id
			ORDER BY last_at DESC
		`, currentUserID, currentUserID)
		if err != nil {
			utils.Error(w, http.StatusInternalServerError, "Could not fetch unread counts")
			return
		}
		defer rows.Close()

		result := make([]UnreadConversation, 0)
		for rows.Next() {
			var uc UnreadConversation
			if err := rows.Scan(&uc.SenderID, &uc.Count, &uc.LastMessage, &uc.LastAt); err != nil {
				continue
			}
			result = append(result, uc)
		}

		utils.JSON(w, http.StatusOK, result)
	}
}

// MarkReadHandler marks all messages from a sender as read.
// POST /messages/mark-read  body: {"sender_id": "..."}
func MarkReadHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			utils.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		var body struct {
			SenderID string `json:"sender_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.SenderID == "" {
			utils.Error(w, http.StatusBadRequest, "sender_id required")
			return
		}

		_, _ = db.Exec(
			`UPDATE private_messages SET is_read = 1
			 WHERE sender_id = ? AND receiver_id = ? AND is_read = 0`,
			body.SenderID, currentUserID,
		)

		utils.JSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func ListGroupMessagesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "group_id required", http.StatusBadRequest)
			return
		}

		ok, err := isAcceptedGroupMemberHandler(db, groupID, currentUserID)
		if err != nil || !ok {
			http.Error(w, "Not a group member", http.StatusForbidden)
			return
		}

		// Optional cursor: load messages sent before this timestamp (for pagination)
		before := r.URL.Query().Get("before")

		var rows *sql.Rows
		if before != "" {
			rows, err = db.Query(`
				SELECT id, group_id, sender_id, content, created_at
				FROM group_messages
				WHERE group_id = ? AND created_at < ?
				ORDER BY created_at DESC
				LIMIT 10
			`, groupID, before)
		} else {
			rows, err = db.Query(`
				SELECT id, group_id, sender_id, content, created_at
				FROM (
					SELECT id, group_id, sender_id, content, created_at
					FROM group_messages
					WHERE group_id = ?
					ORDER BY created_at DESC
					LIMIT 10
				) sub
				ORDER BY created_at ASC
			`, groupID)
		}
		if err != nil {
			http.Error(w, "Could not fetch group messages", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		messages := make([]GroupMessageResponse, 0)

		for rows.Next() {
			var m GroupMessageResponse
			err := rows.Scan(&m.ID, &m.GroupID, &m.SenderID, &m.Content, &m.CreatedAt)
			if err != nil {
				http.Error(w, "Could not read group messages", http.StatusInternalServerError)
				return
			}
			messages = append(messages, m)
		}

		// When paginating backwards, the DB returned DESC; reverse to chronological order
		if before != "" {
			for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
				messages[i], messages[j] = messages[j], messages[i]
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}

func canUsersChatHandlers(db *sql.DB, userA, userB string) bool {
	var count int

	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM followers
		WHERE (
			(follower_id = ? AND following_id = ? AND status = 'accepted')
			OR
			(follower_id = ? AND following_id = ? AND status = 'accepted')
		)
	`, userA, userB, userB, userA).Scan(&count)
	if err == nil && count > 0 {
		return true
	}

	var profileType string
	err = db.QueryRow(`SELECT profile_type FROM users WHERE id = ?`, userB).Scan(&profileType)
	return err == nil && profileType == "public"
}

func isAcceptedGroupMemberHandler(db *sql.DB, groupID, userID string) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM group_members
		WHERE group_id = ? AND user_id = ? AND status = 'accepted'
	`, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

type RecentConversationResponse struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Avatar          string `json:"avatar"`
	Nickname        string `json:"nickname"`
	LastMessage     string `json:"last_message"`
	LastMessageTime string `json:"last_message_time"`
}

func ListRecentConversationsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT 
				u.id, 
				u.first_name, 
				u.last_name, 
				COALESCE(u.avatar, '') as avatar, 
				COALESCE(u.nickname, '') as nickname,
				m.content,
				m.created_at
			FROM (
				SELECT 
					CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END as peer_id,
					content,
					created_at,
					ROW_NUMBER() OVER (
						PARTITION BY CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END 
						ORDER BY created_at DESC
					) as rn
				FROM private_messages
				WHERE sender_id = ? OR receiver_id = ?
			) m
			JOIN users u ON u.id = m.peer_id
			WHERE m.rn = 1
			ORDER BY m.created_at DESC
		`, currentUserID, currentUserID, currentUserID, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch recent conversations", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		conversations := make([]RecentConversationResponse, 0)
		for rows.Next() {
			var c RecentConversationResponse
			err := rows.Scan(
				&c.ID,
				&c.FirstName,
				&c.LastName,
				&c.Avatar,
				&c.Nickname,
				&c.LastMessage,
				&c.LastMessageTime,
			)
			if err != nil {
				http.Error(w, "Could not read recent conversations", http.StatusInternalServerError)
				return
			}
			conversations = append(conversations, c)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(conversations)
	}
}
