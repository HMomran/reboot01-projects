package websocket

import (
	"database/sql"
	"net/http"
)

func PresenceHandler(db *sql.DB, hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := GetCurrentUserIDFromRequest(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		hub.SetPresenceClient(userID, conn)
		defer hub.RemovePresenceClient(userID)

		// Keep alive: read and discard any pings; exit on close
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}
}
