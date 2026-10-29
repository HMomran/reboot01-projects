package auth

import (
	"forum/backend/database"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	DB.DeleteSession(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
