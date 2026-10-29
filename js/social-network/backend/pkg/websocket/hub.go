package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	privateClients  map[string]*websocket.Conn
	presenceClients map[string]*websocket.Conn
	notifClients    map[string]*websocket.Conn
	mu              sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		privateClients:  make(map[string]*websocket.Conn),
		presenceClients: make(map[string]*websocket.Conn),
		notifClients:    make(map[string]*websocket.Conn),
	}
}

func (h *Hub) SetPrivateClient(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.privateClients[userID] = conn
}

func (h *Hub) RemovePrivateClient(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.privateClients, userID)
}

func (h *Hub) GetPrivateClient(userID string) (*websocket.Conn, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conn, ok := h.privateClients[userID]
	return conn, ok
}

func (h *Hub) SetPresenceClient(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.presenceClients[userID] = conn
}

func (h *Hub) RemovePresenceClient(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.presenceClients, userID)
}

func (h *Hub) SetNotifClient(userID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.notifClients[userID] = conn
}

func (h *Hub) RemoveNotifClient(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.notifClients, userID)
}

// PushNotification sends a live notification to a specific user if they are connected.
func (h *Hub) PushNotification(userID string, payload map[string]interface{}) {
	h.mu.RLock()
	conn, ok := h.notifClients[userID]
	h.mu.RUnlock()
	if ok {
		_ = conn.WriteJSON(payload)
	}
}

func (h *Hub) GetOnlineUserIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	seen := make(map[string]struct{})
	for id := range h.presenceClients {
		seen[id] = struct{}{}
	}
	for id := range h.privateClients {
		seen[id] = struct{}{}
	}
	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids
}
