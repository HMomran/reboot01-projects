package websocket

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
)

func GetCurrentUserIDFromRequest(db *sql.DB, r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return "", errors.New("missing session cookie")
	}

	var userID string
	var expiresAt time.Time

	err = db.QueryRow(
		"SELECT user_id, expires_at FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&userID, &expiresAt)
	if err != nil {
		return "", errors.New("invalid session")
	}

	if time.Now().After(expiresAt) {
		return "", errors.New("session expired")
	}

	return userID, nil
}
