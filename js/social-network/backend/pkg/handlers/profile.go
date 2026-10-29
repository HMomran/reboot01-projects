package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social-network/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

type ProfileResponse struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DateOfBirth    string `json:"date_of_birth"`
	Avatar         string `json:"avatar"`
	Nickname       string `json:"nickname"`
	AboutMe        string `json:"about_me"`
	ProfileType    string `json:"profile_type"`
	ShowEmail      bool   `json:"show_email"`
	ShowDOB        bool   `json:"show_dob"`
	IsOwner        bool   `json:"is_owner"`
	IsFollower     bool   `json:"is_follower"`
	FollowStatus   string `json:"follow_status"` // "none" | "pending" | "accepted"
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
	PostsCount     int    `json:"posts_count"`
	CanViewPosts   bool   `json:"can_view_posts"`
}

type UpdatePrivacyRequest struct {
	ProfileType string `json:"profile_type"`
}

func GetProfileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		targetUserID := r.URL.Query().Get("user_id")
		if targetUserID == "" {
			http.Error(w, "Missing user_id", http.StatusBadRequest)
			return
		}

		currentUserID, _ := auth.GetCurrentUserID(db, r)

		var profile ProfileResponse
		var avatar sql.NullString
		var nickname sql.NullString
		var aboutMe sql.NullString
		var showEmail, showDOB int

		err := db.QueryRow(`
			SELECT id, email, first_name, last_name, date_of_birth, avatar, nickname, about_me, profile_type,
			       COALESCE(show_email, 0), COALESCE(show_dob, 0)
			FROM users
			WHERE id = ?
		`, targetUserID).Scan(
			&profile.ID,
			&profile.Email,
			&profile.FirstName,
			&profile.LastName,
			&profile.DateOfBirth,
			&avatar,
			&nickname,
			&aboutMe,
			&profile.ProfileType,
			&showEmail,
			&showDOB,
		)

		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		profile.ShowEmail = showEmail == 1
		profile.ShowDOB = showDOB == 1

		if avatar.Valid { profile.Avatar = avatar.String }
		if nickname.Valid { profile.Nickname = nickname.String }
		if aboutMe.Valid { profile.AboutMe = aboutMe.String }

		profile.IsOwner = currentUserID == targetUserID
		profile.IsFollower = false
		profile.FollowStatus = "none"

		if !profile.IsOwner && currentUserID != "" {
			var followStatus sql.NullString
			err = db.QueryRow(`
				SELECT status FROM followers
				WHERE follower_id = ? AND following_id = ?
				ORDER BY created_at DESC LIMIT 1
			`, currentUserID, targetUserID).Scan(&followStatus)
			if err == nil && followStatus.Valid {
				profile.FollowStatus = followStatus.String
				if followStatus.String == "accepted" {
					profile.IsFollower = true
				}
			}
		}

		// Count followers, following, posts
		_ = db.QueryRow(`SELECT COUNT(*) FROM followers WHERE following_id = ? AND status = 'accepted'`, targetUserID).Scan(&profile.FollowersCount)
		_ = db.QueryRow(`SELECT COUNT(*) FROM followers WHERE follower_id = ? AND status = 'accepted'`, targetUserID).Scan(&profile.FollowingCount)
		_ = db.QueryRow(`SELECT COUNT(*) FROM posts WHERE user_id = ?`, targetUserID).Scan(&profile.PostsCount)

		// Determine if requester can view posts
		profile.CanViewPosts = profile.IsOwner || profile.IsFollower || profile.ProfileType == "public"

		// Hide email/dob from visitors unless the owner opted to show them
		if !profile.IsOwner {
			if !profile.ShowEmail { profile.Email = "" }
			if !profile.ShowDOB  { profile.DateOfBirth = "" }
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}
}

func UpdateProfilePrivacyHandler(db *sql.DB) http.HandlerFunc {
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

		var body UpdatePrivacyRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.ProfileType != "public" && body.ProfileType != "private" {
			http.Error(w, "profile_type must be 'public' or 'private'", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(
			"UPDATE users SET profile_type = ? WHERE id = ?",
			body.ProfileType,
			currentUserID,
		)
		if err != nil {
			http.Error(w, "Could not update profile privacy", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Profile privacy updated successfully"))
	}
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func ChangePasswordHandler(db *sql.DB) http.HandlerFunc {
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
		var body ChangePasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if len(body.NewPassword) < 8 {
			http.Error(w, "New password must be at least 8 characters", http.StatusBadRequest)
			return
		}
		var hashedPassword string
		err = db.QueryRow("SELECT password FROM users WHERE id = ?", currentUserID).Scan(&hashedPassword)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(body.CurrentPassword)); err != nil {
			http.Error(w, "Incorrect current password", http.StatusUnauthorized)
			return
		}
		newHashed, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Could not hash password", http.StatusInternalServerError)
			return
		}
		if _, err = db.Exec("UPDATE users SET password = ? WHERE id = ?", string(newHashed), currentUserID); err != nil {
			http.Error(w, "Could not update password", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password updated successfully"))
	}
}
