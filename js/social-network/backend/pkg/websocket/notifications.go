package websocket

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// NotificationsHandler upgrades to a WebSocket and registers the connection
// in the hub so that real-time notification pushes can reach this user.
func NotificationsHandler(db *sql.DB, hub *Hub) http.HandlerFunc {
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

		hub.SetNotifClient(currentUserID, conn)
		defer hub.RemoveNotifClient(currentUserID)

		// Keep connection alive with periodic pings.
		conn.SetReadDeadline(time.Now().Add(70 * time.Second))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(70 * time.Second))
			return nil
		})

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				// Drain any client-sent messages (keeps the read loop alive)
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}
