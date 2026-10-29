package middleware

import (
	"database/sql"
	"net/http"
	"time"
)

func AuthMiddleware(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var expiresAt time.Time
		err = db.QueryRow(
			"SELECT expires_at FROM sessions WHERE id = ?",
			cookie.Value,
		).Scan(&expiresAt)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		if time.Now().After(expiresAt) {
			http.Error(w, "Session expired", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
