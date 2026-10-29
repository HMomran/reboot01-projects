package websocket

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GroupIncomingMessage struct {
	GroupID string `json:"group_id"`
	Content string `json:"content"`
}

type GroupOutgoingMessage struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type GroupHub struct {
	groupClients map[string]map[string]*websocket.Conn
	mu           sync.RWMutex
	// NotifHub is used to push group-message notifications to members
	// who are not currently connected to this group's WS.
	NotifHub interface {
		PushNotification(userID string, payload map[string]interface{})
	}
}

func NewGroupHub() *GroupHub {
	return &GroupHub{
		groupClients: make(map[string]map[string]*websocket.Conn),
	}
}

func (g *GroupHub) AddClient(groupID, userID string, conn *websocket.Conn) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.groupClients[groupID]; !exists {
		g.groupClients[groupID] = make(map[string]*websocket.Conn)
	}

	g.groupClients[groupID][userID] = conn
}

func (g *GroupHub) RemoveClient(groupID, userID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.groupClients[groupID]; exists {
		delete(g.groupClients[groupID], userID)
		if len(g.groupClients[groupID]) == 0 {
			delete(g.groupClients, groupID)
		}
	}
}

func (g *GroupHub) Broadcast(groupID string, msg GroupOutgoingMessage) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if clients, exists := g.groupClients[groupID]; exists {
		for _, conn := range clients {
			_ = conn.WriteJSON(msg)
		}
	}
}

func GroupChatHandler(db *sql.DB, groupHub *GroupHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := GetCurrentUserIDFromRequest(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "group_id required", http.StatusBadRequest)
			return
		}

		if !isAcceptedGroupMember(db, groupID, currentUserID) {
			http.Error(w, "Not a group member", http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		groupHub.AddClient(groupID, currentUserID, conn)
		defer groupHub.RemoveClient(groupID, currentUserID)

		for {
			var msg GroupIncomingMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}

			if msg.Content == "" || len(msg.Content) > 1000 {
				continue
			}

			messageID := uuid.New().String()

			_, err = db.Exec(`
				INSERT INTO group_messages (id, group_id, sender_id, content)
				VALUES (?, ?, ?, ?)
			`, messageID, groupID, currentUserID, msg.Content)
			if err != nil {
				continue
			}

			var createdAt string
			_ = db.QueryRow(`
				SELECT created_at
				FROM group_messages
				WHERE id = ?
			`, messageID).Scan(&createdAt)

			out := GroupOutgoingMessage{
				ID:        messageID,
				GroupID:   groupID,
				SenderID:  currentUserID,
				Content:   msg.Content,
				CreatedAt: createdAt,
			}

			groupHub.Broadcast(groupID, out)

			// Push a live notification to every group member NOT currently
			// connected to this group's WS (i.e. not actively viewing it).
			if groupHub.NotifHub != nil {
				// Fetch group title for the notification message
				var groupTitle string
				_ = db.QueryRow(`SELECT title FROM groups WHERE id = ?`, groupID).Scan(&groupTitle)

				// Fetch sender display name
				var senderFirst, senderLast string
				_ = db.QueryRow(`SELECT first_name, last_name FROM users WHERE id = ?`, currentUserID).Scan(&senderFirst, &senderLast)

				senderName := senderFirst + " " + senderLast
				notifMsg := fmt.Sprintf("%s in %s: %s", senderName, groupTitle, truncate(msg.Content, 60))

				// Get all accepted members except the sender
				rows2, err2 := db.Query(`
					SELECT user_id FROM group_members
					WHERE group_id = ? AND user_id != ? AND status = 'accepted'
				`, groupID, currentUserID)
				if err2 == nil {
					defer rows2.Close()
					groupHub.mu.RLock()
					activeInGroup := groupHub.groupClients[groupID]
					groupHub.mu.RUnlock()
					for rows2.Next() {
						var memberID string
						if rows2.Scan(&memberID) == nil {
							// Only notify if not currently viewing this group chat
							if _, inGroup := activeInGroup[memberID]; !inGroup {
								groupHub.NotifHub.PushNotification(memberID, map[string]interface{}{
									"event":    "notification",
									"type":     "group_message",
									"message":  notifMsg,
									"group_id": groupID,
								})
							}
						}
					}
				}
			}
		}
	}
}

func isAcceptedGroupMember(db *sql.DB, groupID, userID string) bool {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM group_members
		WHERE group_id = ? AND user_id = ? AND status = 'accepted'
	`, groupID, userID).Scan(&count)
	return err == nil && count > 0
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "…"
}
