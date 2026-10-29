package websocket

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// FeedHub broadcasts post-update events to every user currently viewing their feed.
type FeedHub struct {
	clients map[string]*websocket.Conn // userID → conn
	mu      sync.RWMutex
}

func NewFeedHub() *FeedHub {
	return &FeedHub{clients: make(map[string]*websocket.Conn)}
}

func (f *FeedHub) Add(userID string, conn *websocket.Conn) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.clients[userID] = conn
}

func (f *FeedHub) Remove(userID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.clients, userID)
}

// Broadcast sends payload to every connected feed client.
func (f *FeedHub) Broadcast(payload map[string]interface{}) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, conn := range f.clients {
		_ = conn.WriteJSON(payload)
	}
}

// FeedHandler upgrades to WS and keeps the connection alive.
// The server pushes "post_updated" events; the client never sends anything.
func FeedHandler(db *sql.DB, hub *FeedHub) http.HandlerFunc {
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

		hub.Add(currentUserID, conn)
		defer hub.Remove(currentUserID)

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
