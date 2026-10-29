package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateSession(db *sql.DB, userID string, w http.ResponseWriter) error {
	sessionID := uuid.New().String()
	expiry := time.Now().Add(24 * time.Hour)

	_, err := db.Exec(
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID,
		userID,
		expiry,
	)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiry,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	return nil
}
