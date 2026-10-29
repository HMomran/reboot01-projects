package auth

import (
	"database/sql"
	"net/http"
)

func LogoutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "No active session", http.StatusUnauthorized)
			return
		}

		_, err = db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
		if err != nil {
			http.Error(w, "Could not logout", http.StatusInternalServerError)
			return
		}

		expiredCookie := &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, expiredCookie)

		w.Write([]byte("Logged out successfully"))
	}
}
