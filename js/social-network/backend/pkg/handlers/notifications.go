package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social-network/pkg/auth"
)

type NotificationResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Type        string `json:"type"`
	ReferenceID string `json:"reference_id"`
	Message     string `json:"message"`
	IsRead      bool   `json:"is_read"`
	CreatedAt   string `json:"created_at"`
}

type MarkNotificationReadRequest struct {
	NotificationID string `json:"notification_id"`
}

func ListNotificationsHandler(db *sql.DB) http.HandlerFunc {
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
			SELECT id, user_id, type, reference_id, message, is_read, created_at
			FROM notifications
			WHERE user_id = ?
			ORDER BY created_at DESC
		`, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch notifications", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		notifications := make([]NotificationResponse, 0)

		for rows.Next() {
			var item NotificationResponse
			var referenceID sql.NullString
			var isReadInt int

			err := rows.Scan(
				&item.ID,
				&item.UserID,
				&item.Type,
				&referenceID,
				&item.Message,
				&isReadInt,
				&item.CreatedAt,
			)
			if err != nil {
				http.Error(w, "Could not read notifications", http.StatusInternalServerError)
				return
			}

			if referenceID.Valid {
				item.ReferenceID = referenceID.String
			} else {
				item.ReferenceID = ""
			}

			item.IsRead = isReadInt == 1
			notifications = append(notifications, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifications)
	}
}

func MarkNotificationReadHandler(db *sql.DB) http.HandlerFunc {
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

		var body MarkNotificationReadRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.NotificationID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(`
			UPDATE notifications
			SET is_read = 1
			WHERE id = ? AND user_id = ?
		`, body.NotificationID, currentUserID)
		if err != nil {
			http.Error(w, "Could not update notification", http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}

		w.Write([]byte("Notification marked as read"))
	}
}

func UnreadNotificationsCountHandler(db *sql.DB) http.HandlerFunc {
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
		var count int
		err = db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0`, currentUserID).Scan(&count)
		if err != nil {
			http.Error(w, "Could not fetch count", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"unread": count})
	}
}

func MarkAllNotificationsReadHandler(db *sql.DB) http.HandlerFunc {
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
		_, err = db.Exec(`UPDATE notifications SET is_read = 1 WHERE user_id = ?`, currentUserID)
		if err != nil {
			http.Error(w, "Could not update notifications", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("All notifications marked as read"))
	}
}

func ClearAllNotificationsHandler(db *sql.DB) http.HandlerFunc {
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
		_, err = db.Exec(`DELETE FROM notifications WHERE user_id = ?`, currentUserID)
		if err != nil {
			http.Error(w, "Could not delete notifications", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("All notifications cleared"))
	}
}
