package websocket

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// PrivateIncomingMessage is what clients send over the WS connection.
// Type can be:
//
//	"message"  – send a chat message (requires receiver_id + content)
//	"typing"   – notify the other user you are typing (requires receiver_id)
//	"read"     – mark messages from sender_id as read (requires sender_id)
type PrivateIncomingMessage struct {
	Type       string `json:"type"`
	ReceiverID string `json:"receiver_id"`
	SenderID   string `json:"sender_id"`
	Content    string `json:"content"`
}

// PrivateOutgoingMessage is what the server pushes to clients.
type PrivateOutgoingMessage struct {
	Type       string `json:"type"`
	ID         string `json:"id,omitempty"`
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name,omitempty"`
	ReceiverID string `json:"receiver_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	IsRead     bool   `json:"is_read,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func PrivateChatHandler(db *sql.DB, hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := GetCurrentUserIDFromRequest(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		hub.SetPrivateClient(currentUserID, conn)
		defer hub.RemovePrivateClient(currentUserID)

		for {
			var msg PrivateIncomingMessage
			if err := conn.ReadJSON(&msg); err != nil {
				break
			}

			switch msg.Type {
			case "typing", "stop_typing":
				// Relay typing / stop_typing notification to the receiver only – no DB write
				if msg.ReceiverID == "" {
					continue
				}
				if receiverConn, ok := hub.GetPrivateClient(msg.ReceiverID); ok {
					_ = receiverConn.WriteJSON(PrivateOutgoingMessage{
						Type:     msg.Type,
						SenderID: currentUserID,
					})
				}

			case "read":
				// Mark all unread messages from msg.SenderID → currentUserID as read
				if msg.SenderID == "" {
					continue
				}
				_, _ = db.Exec(
					`UPDATE private_messages SET is_read = 1
					 WHERE sender_id = ? AND receiver_id = ? AND is_read = 0`,
					msg.SenderID, currentUserID,
				)
				// Notify the original sender that their messages were seen
				if senderConn, ok := hub.GetPrivateClient(msg.SenderID); ok {
					_ = senderConn.WriteJSON(PrivateOutgoingMessage{
						Type:     "read",
						SenderID: currentUserID, // "currentUser has read your messages"
					})
				}

			default:
				// "message" or legacy (empty type) – save and relay
				if msg.ReceiverID == "" || msg.Content == "" {
					continue
				}
				if len(msg.Content) > 1000 {
					continue
				}
				if !canUsersChat(db, currentUserID, msg.ReceiverID) {
					continue
				}

				messageID := uuid.New().String()
				_, err = db.Exec(`
					INSERT INTO private_messages (id, sender_id, receiver_id, content, is_read)
					VALUES (?, ?, ?, ?, 0)
				`, messageID, currentUserID, msg.ReceiverID, msg.Content)
				if err != nil {
					continue
				}

				var createdAt string
				_ = db.QueryRow(
					`SELECT created_at FROM private_messages WHERE id = ?`, messageID,
				).Scan(&createdAt)

				var senderFirst, senderLast string
				_ = db.QueryRow("SELECT first_name, last_name FROM users WHERE id = ?", currentUserID).Scan(&senderFirst, &senderLast)
				senderName := senderFirst + " " + senderLast

				out := PrivateOutgoingMessage{
					Type:       "message",
					ID:         messageID,
					SenderID:   currentUserID,
					SenderName: senderName,
					ReceiverID: msg.ReceiverID,
					Content:    msg.Content,
					CreatedAt:  createdAt,
					IsRead:     false,
				}

				// Echo back to sender
				_ = conn.WriteJSON(out)

				// Deliver to receiver if online
				if receiverConn, ok := hub.GetPrivateClient(msg.ReceiverID); ok {
					_ = receiverConn.WriteJSON(out)
				}
			}
		}
	}
}

func canUsersChat(db *sql.DB, userA, userB string) bool {
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
	if err != nil {
		return false
	}

	return count > 0
}
