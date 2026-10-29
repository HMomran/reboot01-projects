package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"social-network/pkg/auth"
)

type OnlineChecker interface {
	GetOnlineUserIDs() []string
}

type MeResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Avatar      string `json:"avatar"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"about_me"`
	DateOfBirth string `json:"date_of_birth"`
	ProfileType string `json:"profile_type"`
	ShowEmail   bool   `json:"show_email"`
	ShowDOB     bool   `json:"show_dob"`
}

type UserListItem struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
}

func OnlineUsersHandler(db *sql.DB, hub OnlineChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, _ := auth.GetCurrentUserID(db, r)
		all := hub.GetOnlineUserIDs()
		filtered := make([]string, 0, len(all))
		for _, id := range all {
			if id != currentUserID {
				filtered = append(filtered, id)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filtered)
	}
}

func ListUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, first_name, last_name, COALESCE(nickname, ''), COALESCE(avatar, '')
			FROM users
			WHERE id != ?
			ORDER BY first_name ASC
		`, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		users := make([]UserListItem, 0)
		for rows.Next() {
			var u UserListItem
			if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Nickname, &u.Avatar); err != nil {
				http.Error(w, "Could not read users", http.StatusInternalServerError)
				return
			}
			users = append(users, u)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func MeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentUserID, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var me MeResponse
		var avatar sql.NullString
		var nickname sql.NullString
		var aboutMe sql.NullString
		var showEmail, showDOB int

		err = db.QueryRow(`
			SELECT id, email, first_name, last_name, avatar, nickname, about_me, date_of_birth, profile_type,
			       COALESCE(show_email, 0), COALESCE(show_dob, 0)
			FROM users
			WHERE id = ?
		`, currentUserID).Scan(
			&me.ID, &me.Email, &me.FirstName, &me.LastName,
			&avatar, &nickname, &aboutMe, &me.DateOfBirth, &me.ProfileType,
			&showEmail, &showDOB,
		)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if avatar.Valid { me.Avatar = avatar.String }
		if nickname.Valid { me.Nickname = nickname.String }
		if aboutMe.Valid { me.AboutMe = aboutMe.String }
		me.ShowEmail = showEmail == 1
		me.ShowDOB = showDOB == 1

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(me)
	}
}

// UpdateProfileHandler allows updating nickname and about_me
func UpdateProfileHandler(db *sql.DB) http.HandlerFunc {
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
		var body struct {
			Nickname  string `json:"nickname"`
			AboutMe   string `json:"about_me"`
			ShowEmail *bool  `json:"show_email"`
			ShowDOB   *bool  `json:"show_dob"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		var nickVal, aboutVal interface{}
		if body.Nickname != "" {
			body.Nickname = strings.TrimSpace(body.Nickname)
			if len(body.Nickname) > 30 {
				http.Error(w, "Nickname cannot exceed 30 characters", http.StatusBadRequest)
				return
			}
			nickFormatRe := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
			if !nickFormatRe.MatchString(body.Nickname) {
				http.Error(w, "Nickname can only contain letters, numbers, and underscores", http.StatusBadRequest)
				return
			}
			var count int
			err = db.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(nickname) = LOWER(?) AND id != ?", body.Nickname, currentUserID).Scan(&count)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			if count > 0 {
				http.Error(w, "Nickname is already taken. Please choose another one.", http.StatusBadRequest)
				return
			}
			nickVal = body.Nickname
		}
		if body.AboutMe != "" {
			aboutVal = body.AboutMe
		}

		// Resolve show_email and show_dob — default to current DB value if not provided
		var curShowEmail, curShowDOB int
		_ = db.QueryRow("SELECT COALESCE(show_email,0), COALESCE(show_dob,0) FROM users WHERE id = ?", currentUserID).Scan(&curShowEmail, &curShowDOB)
		showEmailVal := curShowEmail
		showDOBVal := curShowDOB
		if body.ShowEmail != nil {
			if *body.ShowEmail { showEmailVal = 1 } else { showEmailVal = 0 }
		}
		if body.ShowDOB != nil {
			if *body.ShowDOB { showDOBVal = 1 } else { showDOBVal = 0 }
		}

		_, err = db.Exec(
			"UPDATE users SET nickname = ?, about_me = ?, show_email = ?, show_dob = ? WHERE id = ?",
			nickVal, aboutVal, showEmailVal, showDOBVal, currentUserID,
		)
		if err != nil {
			http.Error(w, "Could not update profile", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Profile updated"))
	}
}
