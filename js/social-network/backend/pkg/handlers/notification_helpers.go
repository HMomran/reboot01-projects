package handlers

import (
	"database/sql"

	"github.com/google/uuid"
)

// NotifPusher is satisfied by *websocket.Hub. Using an interface here avoids
// an import cycle between the handlers and websocket packages.
type NotifPusher interface {
	PushNotification(userID string, payload map[string]interface{})
}

// FeedBroadcaster is satisfied by *websocket.FeedHub.
type FeedBroadcaster interface {
	Broadcast(payload map[string]interface{})
}

// GlobalNotifHub is set by main.go after creating the Hub.
// When non-nil, every new notification is also pushed live to the recipient.
var GlobalNotifHub NotifPusher

// GlobalFeedHub is set by main.go. When non-nil, post-update events are
// pushed to all connected feed clients so they see live reactions/comments.
var GlobalFeedHub FeedBroadcaster

func createNotification(db *sql.DB, userID, notifType, referenceID, message string) error {
	id := uuid.New().String()
	_, err := db.Exec(`
		INSERT INTO notifications (id, user_id, type, reference_id, message)
		VALUES (?, ?, ?, ?, ?)
	`,
		id,
		userID,
		notifType,
		referenceID,
		message,
	)

	if err != nil {
		return err
	}

	// Push live to the recipient if they are connected via /ws/notifications
	if GlobalNotifHub != nil {
		GlobalNotifHub.PushNotification(userID, map[string]interface{}{
			"event":        "notification",
			"id":           id,
			"type":         notifType,
			"reference_id": referenceID,
			"message":      message,
		})
	}

	return nil
}
