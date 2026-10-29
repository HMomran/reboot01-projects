package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social-network/pkg/auth"

	"github.com/google/uuid"
)

type FollowRequestBody struct {
	TargetUserID string `json:"target_user_id"`
}

type FollowActionBody struct {
	RequestID string `json:"request_id"`
}

func FollowUserHandler(db *sql.DB) http.HandlerFunc {
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

		var body FollowRequestBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.TargetUserID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.TargetUserID == currentUserID {
			http.Error(w, "You cannot follow yourself", http.StatusBadRequest)
			return
		}

		var profileType string
		err = db.QueryRow(
			"SELECT profile_type FROM users WHERE id = ?",
			body.TargetUserID,
		).Scan(&profileType)
		if err != nil {
			http.Error(w, "Target user not found", http.StatusNotFound)
			return
		}

		status := "pending"
		if profileType == "public" {
			status = "accepted"
		}

		followID := uuid.New().String()

		_, err = db.Exec(
			"INSERT INTO followers (id, follower_id, following_id, status) VALUES (?, ?, ?, ?)",
			followID,
			currentUserID,
			body.TargetUserID,
			status,
		)
		if err != nil {
			http.Error(w, "Follow request already exists or could not be created", http.StatusBadRequest)
			return
		}

		var firstName string
		var lastName string
		err = db.QueryRow(
			"SELECT first_name, last_name FROM users WHERE id = ?",
			currentUserID,
		).Scan(&firstName, &lastName)

		if err == nil {
			switch status {
			case "pending":
				_ = createNotification(
					db,
					body.TargetUserID,
					"follow_request",
					currentUserID,
					firstName+" "+lastName+" sent you a follow request",
				)
			case "accepted":
				_ = createNotification(
					db,
					body.TargetUserID,
					"new_follower",
					currentUserID,
					firstName+" "+lastName+" started following you",
				)
			}
		}

		if status == "accepted" {
			w.Write([]byte("User followed successfully"))
			return
		}

		w.Write([]byte("Follow request sent"))
	}
}

func AcceptFollowRequestHandler(db *sql.DB) http.HandlerFunc {
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

		var body FollowActionBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.RequestID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// reference_id in the notification is the follower's user ID.
		// Accept by matching follower_id (= reference_id sent as request_id).
		followerID := body.RequestID

		result, err := db.Exec(
			"UPDATE followers SET status = 'accepted' WHERE follower_id = ? AND following_id = ? AND status = 'pending'",
			followerID,
			currentUserID,
		)
		if err != nil {
			http.Error(w, "Could not accept request", http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "Follow request not found", http.StatusNotFound)
			return
		}

		// Create acceptance notification for the follower
		var firstName, lastName string
		_ = db.QueryRow("SELECT first_name, last_name FROM users WHERE id = ?", currentUserID).Scan(&firstName, &lastName)
		_ = createNotification(
			db,
			followerID,
			"follow_accept",
			currentUserID,
			firstName+" "+lastName+" accepted your follow request",
		)

		w.Write([]byte("Follow request accepted"))
	}
}

func DeclineFollowRequestHandler(db *sql.DB) http.HandlerFunc {
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

		var body FollowActionBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.RequestID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"UPDATE followers SET status = 'declined' WHERE follower_id = ? AND following_id = ? AND status = 'pending'",
			body.RequestID,
			currentUserID,
		)
		if err != nil {
			http.Error(w, "Could not decline request", http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "Follow request not found", http.StatusNotFound)
			return
		}

		w.Write([]byte("Follow request declined"))
	}
}

func UnfollowUserHandler(db *sql.DB) http.HandlerFunc {
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

		var body FollowRequestBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.TargetUserID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"DELETE FROM followers WHERE follower_id = ? AND following_id = ? AND (status = 'accepted' OR status = 'pending')",
			currentUserID,
			body.TargetUserID,
		)
		if err != nil {
			http.Error(w, "Could not unfollow user", http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			http.Error(w, "Follow relationship not found", http.StatusNotFound)
			return
		}

		w.Write([]byte("User unfollowed successfully"))
	}
}

func ListPendingFollowRequestsHandler(db *sql.DB) http.HandlerFunc {
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
			SELECT f.id, u.id, u.first_name, u.last_name, u.nickname
			FROM followers f
			JOIN users u ON f.follower_id = u.id
			WHERE f.following_id = ? AND f.status = 'pending'
			ORDER BY f.created_at DESC
		`, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch pending requests", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type PendingRequest struct {
			RequestID string `json:"request_id"`
			UserID    string `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Nickname  string `json:"nickname"`
		}

		var requests []PendingRequest

		for rows.Next() {
			var item PendingRequest
			var nickname sql.NullString

			err := rows.Scan(&item.RequestID, &item.UserID, &item.FirstName, &item.LastName, &nickname)
			if err != nil {
				http.Error(w, "Could not read pending requests", http.StatusInternalServerError)
				return
			}
			if nickname.Valid {
				item.Nickname = nickname.String
			}
			requests = append(requests, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requests)
	}
}

func ListFollowersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "Missing user_id", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
			SELECT u.id, u.first_name, u.last_name, u.nickname
			FROM followers f
			JOIN users u ON f.follower_id = u.id
			WHERE f.following_id = ? AND f.status = 'accepted'
			ORDER BY f.created_at DESC
		`, userID)
		if err != nil {
			http.Error(w, "Could not fetch followers", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type UserSummary struct {
			UserID    string `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Nickname  string `json:"nickname"`
		}

		followers := make([]UserSummary, 0)

		for rows.Next() {
			var item UserSummary
			var nickname sql.NullString

			err := rows.Scan(&item.UserID, &item.FirstName, &item.LastName, &nickname)
			if err != nil {
				http.Error(w, "Could not read followers", http.StatusInternalServerError)
				return
			}

			if nickname.Valid {
				item.Nickname = nickname.String
			} else {
				item.Nickname = ""
			}

			followers = append(followers, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(followers)
	}
}

func ListFollowingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "Missing user_id", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
			SELECT u.id, u.first_name, u.last_name, u.nickname
			FROM followers f
			JOIN users u ON f.following_id = u.id
			WHERE f.follower_id = ? AND f.status = 'accepted'
			ORDER BY f.created_at DESC
		`, userID)
		if err != nil {
			http.Error(w, "Could not fetch following users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type UserSummary struct {
			UserID    string `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Nickname  string `json:"nickname"`
		}

		following := make([]UserSummary, 0)

		for rows.Next() {
			var item UserSummary
			var nickname sql.NullString

			err := rows.Scan(&item.UserID, &item.FirstName, &item.LastName, &nickname)
			if err != nil {
				http.Error(w, "Could not read following users", http.StatusInternalServerError)
				return
			}

			if nickname.Valid {
				item.Nickname = nickname.String
			} else {
				item.Nickname = ""
			}

			following = append(following, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(following)
	}
}
