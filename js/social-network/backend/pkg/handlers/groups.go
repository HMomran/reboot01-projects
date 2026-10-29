package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"social-network/pkg/auth"

	"github.com/google/uuid"
)

type GroupMemberResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	IsCreator bool   `json:"is_creator"`
}

func GroupMembersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		_, err := auth.GetCurrentUserID(db, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "group_id required", http.StatusBadRequest)
			return
		}
		rows, err := db.Query(`
			SELECT u.id, u.first_name, u.last_name, COALESCE(u.avatar, '') as avatar,
			       CASE WHEN g.creator_id = u.id THEN 1 ELSE 0 END as is_creator
			FROM group_members gm
			JOIN users u ON u.id = gm.user_id
			JOIN groups g ON g.id = gm.group_id
			WHERE gm.group_id = ? AND gm.status = 'accepted'
		`, groupID)
		if err != nil {
			http.Error(w, "Could not fetch members", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		members := make([]GroupMemberResponse, 0)
		for rows.Next() {
			var m GroupMemberResponse
			if err := rows.Scan(&m.ID, &m.FirstName, &m.LastName, &m.Avatar, &m.IsCreator); err != nil {
				continue
			}
			members = append(members, m)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(members)
	}
}

func deleteGroupCascade(db *sql.DB, groupID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Delete notifications related to group events
	if _, err := tx.Exec(`DELETE FROM notifications WHERE reference_id IN (SELECT id FROM events WHERE group_id = ?)`, groupID); err != nil {
		return err
	}
	// 2. Delete notifications related to group posts
	if _, err := tx.Exec(`DELETE FROM notifications WHERE reference_id IN (SELECT id FROM group_posts WHERE group_id = ?)`, groupID); err != nil {
		return err
	}
	// 3. Delete notifications related to the group itself
	if _, err := tx.Exec(`DELETE FROM notifications WHERE reference_id = ?`, groupID); err != nil {
		return err
	}
	// 4. Delete event responses
	if _, err := tx.Exec(`DELETE FROM event_responses WHERE event_id IN (SELECT id FROM events WHERE group_id = ?)`, groupID); err != nil {
		return err
	}
	// 5. Delete events
	if _, err := tx.Exec(`DELETE FROM events WHERE group_id = ?`, groupID); err != nil {
		return err
	}
	// 6. Delete group comments
	if _, err := tx.Exec(`DELETE FROM group_comments WHERE post_id IN (SELECT id FROM group_posts WHERE group_id = ?)`, groupID); err != nil {
		return err
	}
	// 7. Delete group posts
	if _, err := tx.Exec(`DELETE FROM group_posts WHERE group_id = ?`, groupID); err != nil {
		return err
	}
	// 8. Delete group messages
	if _, err := tx.Exec(`DELETE FROM group_messages WHERE group_id = ?`, groupID); err != nil {
		return err
	}
	// 9. Delete group members
	if _, err := tx.Exec(`DELETE FROM group_members WHERE group_id = ?`, groupID); err != nil {
		return err
	}
	// 10. Delete group itself
	if _, err := tx.Exec(`DELETE FROM groups WHERE id = ?`, groupID); err != nil {
		return err
	}

	return tx.Commit()
}

func LeaveGroupHandler(db *sql.DB) http.HandlerFunc {
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
		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "group_id required", http.StatusBadRequest)
			return
		}
		var creatorID string
		err = db.QueryRow(`SELECT creator_id FROM groups WHERE id = ?`, groupID).Scan(&creatorID)
		if err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}
		if creatorID == currentUserID {
			// Transfer admin to another member if one exists
			var newCreatorID string
			err = db.QueryRow(
				`SELECT user_id FROM group_members WHERE group_id = ? AND user_id != ? AND status = 'accepted' LIMIT 1`,
				groupID, currentUserID,
			).Scan(&newCreatorID)
			if err != nil {
				// No other members — dissolve the group
				if err := deleteGroupCascade(db, groupID); err != nil {
					http.Error(w, "Could not dissolve group: "+err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}
			_, err = db.Exec(`UPDATE groups SET creator_id = ? WHERE id = ?`, newCreatorID, groupID)
			if err != nil {
				http.Error(w, "Could not transfer ownership: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
		_, err = db.Exec(`DELETE FROM group_members WHERE group_id = ? AND user_id = ?`, groupID, currentUserID)
		if err != nil {
			http.Error(w, "Could not leave group: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteGroupHandler(db *sql.DB) http.HandlerFunc {
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
		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "group_id required", http.StatusBadRequest)
			return
		}
		// Only creator can delete
		var creatorID string
		err = db.QueryRow(`SELECT creator_id FROM groups WHERE id = ?`, groupID).Scan(&creatorID)
		if err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}
		if creatorID != currentUserID {
			http.Error(w, "Only the group creator can delete this group", http.StatusForbidden)
			return
		}
		if err := deleteGroupCascade(db, groupID); err != nil {
			http.Error(w, "Could not delete group: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type CreateGroupRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Invitees    []string `json:"invitees"`
}

type GroupResponse struct {
	ID               string `json:"id"`
	CreatorID        string `json:"creator_id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	CreatedAt        string `json:"created_at"`
	MemberCount      int    `json:"member_count"`
	MembershipStatus string `json:"membership_status"`
}

type GroupActionRequest struct {
	GroupID string `json:"group_id"`
	UserID  string `json:"user_id"`
}

func CreateGroupHandler(db *sql.DB) http.HandlerFunc {
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

		var body CreateGroupRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.Title == "" || body.Description == "" {
			http.Error(w, "title and description are required", http.StatusBadRequest)
			return
		}
		body.Title = strings.TrimSpace(body.Title)
		body.Description = strings.TrimSpace(body.Description)
		if len(body.Title) > 100 {
			http.Error(w, "Group title cannot exceed 100 characters", http.StatusBadRequest)
			return
		}
		if len(body.Description) > 1000 {
			http.Error(w, "Group description cannot exceed 1000 characters", http.StatusBadRequest)
			return
		}

		groupID := uuid.New().String()

		_, err = db.Exec(`
			INSERT INTO groups (id, creator_id, title, description)
			VALUES (?, ?, ?, ?)
		`, groupID, currentUserID, body.Title, body.Description)
		if err != nil {
			http.Error(w, "Could not create group", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			INSERT INTO group_members (id, group_id, user_id, status)
			VALUES (?, ?, ?, 'accepted')
		`, uuid.New().String(), groupID, currentUserID)
		if err != nil {
			http.Error(w, "Could not add creator to group", http.StatusInternalServerError)
			return
		}

		var creatorFirstName string
		var creatorLastName string
		_ = db.QueryRow(
			"SELECT first_name, last_name FROM users WHERE id = ?",
			currentUserID,
		).Scan(&creatorFirstName, &creatorLastName)

		for _, inviteeID := range body.Invitees {
			if inviteeID == "" || inviteeID == currentUserID {
				continue
			}

			_, err = db.Exec(`
				INSERT OR IGNORE INTO group_members (id, group_id, user_id, status)
				VALUES (?, ?, ?, 'invited')
			`, uuid.New().String(), groupID, inviteeID)
			if err == nil {
				_ = createNotification(
					db,
					inviteeID,
					"group_invite",
					groupID,
					creatorFirstName+" "+creatorLastName+" invited you to join group: "+body.Title,
				)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GroupResponse{
			ID:    groupID,
			Title: body.Title,
		})
	}
}

func ListGroupsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		currentUserID, _ := auth.GetCurrentUserID(db, r)

		rows, err := db.Query(`
			SELECT g.id, g.creator_id, g.title, g.description, g.created_at,
			       (SELECT COUNT(*) FROM group_members WHERE group_id = g.id AND status = 'accepted') AS member_count,
			       COALESCE((SELECT status FROM group_members WHERE group_id = g.id AND user_id = ?), '') AS membership_status
			FROM groups g
			ORDER BY g.created_at DESC
		`, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch groups", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		groups := make([]GroupResponse, 0)

		for rows.Next() {
			var item GroupResponse
			err := rows.Scan(
				&item.ID,
				&item.CreatorID,
				&item.Title,
				&item.Description,
				&item.CreatedAt,
				&item.MemberCount,
				&item.MembershipStatus,
			)
			if err != nil {
				http.Error(w, "Could not read groups", http.StatusInternalServerError)
				return
			}
			groups = append(groups, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}

func RequestToJoinGroupHandler(db *sql.DB) http.HandlerFunc {
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

		var body GroupActionRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.GroupID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var creatorID string
		var groupTitle string

		err = db.QueryRow(`
			SELECT creator_id, title
			FROM groups
			WHERE id = ?
		`, body.GroupID).Scan(&creatorID, &groupTitle)
		if err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		_, err = db.Exec(`
			INSERT INTO group_members (id, group_id, user_id, status)
			VALUES (?, ?, ?, 'requested')
		`, uuid.New().String(), body.GroupID, currentUserID)
		if err != nil {
			http.Error(w, "Request already exists or could not be created", http.StatusBadRequest)
			return
		}

		var requesterFirstName string
		var requesterLastName string
		_ = db.QueryRow(
			"SELECT first_name, last_name FROM users WHERE id = ?",
			currentUserID,
		).Scan(&requesterFirstName, &requesterLastName)

		_ = createNotification(
			db,
			creatorID,
			"group_join_request",
			currentUserID+":"+body.GroupID,
			requesterFirstName+" "+requesterLastName+" requested to join group: "+groupTitle,
		)

		w.Write([]byte("Join request sent"))
	}
}

func InviteToGroupHandler(db *sql.DB) http.HandlerFunc {
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

		var body GroupActionRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil || body.GroupID == "" || body.UserID == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.UserID == currentUserID {
			http.Error(w, "You cannot invite yourself", http.StatusBadRequest)
			return
		}

		isMember, err := isAcceptedGroupMember(db, body.GroupID, currentUserID)
		if err != nil || !isMember {
			http.Error(w, "Only group members can invite users", http.StatusForbidden)
			return
		}

		var groupTitle string
		err = db.QueryRow("SELECT title FROM groups WHERE id = ?", body.GroupID).Scan(&groupTitle)
		if err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		_, err = db.Exec(`
			INSERT INTO group_members (id, group_id, user_id, status)
			VALUES (?, ?, ?, 'invited')
		`, uuid.New().String(), body.GroupID, body.UserID)
		if err != nil {
			http.Error(w, "User is already a member or could not be added", http.StatusBadRequest)
			return
		}

		var inviterFirstName string
		var inviterLastName string
		_ = db.QueryRow(
			"SELECT first_name, last_name FROM users WHERE id = ?",
			currentUserID,
		).Scan(&inviterFirstName, &inviterLastName)

		_ = createNotification(
			db,
			body.UserID,
			"group_invite",
			body.GroupID,
			inviterFirstName+" "+inviterLastName+" invited you to join group: "+groupTitle,
		)

		w.Write([]byte("Invitation sent successfully"))
	}
}

func AcceptGroupInviteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleGroupMembershipStatusUpdate(db, w, r, "invited", "accepted", "Group invite accepted")
	}
}

func DeclineGroupInviteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleGroupMembershipStatusUpdate(db, w, r, "invited", "declined", "Group invite declined")
	}
}

func AcceptGroupJoinRequestHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleCreatorGroupRequestStatusUpdate(db, w, r, "requested", "accepted", "Join request accepted")
	}
}

func DeclineGroupJoinRequestHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleCreatorGroupRequestStatusUpdate(db, w, r, "requested", "declined", "Join request declined")
	}
}

func ListMyGroupInvitesHandler(db *sql.DB) http.HandlerFunc {
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
			SELECT g.id, g.creator_id, g.title, g.description, g.created_at
			FROM group_members gm
			JOIN groups g ON gm.group_id = g.id
			WHERE gm.user_id = ? AND gm.status = 'invited'
			ORDER BY gm.created_at DESC
		`, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch invites", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		groups := make([]GroupResponse, 0)

		for rows.Next() {
			var item GroupResponse
			err := rows.Scan(&item.ID, &item.CreatorID, &item.Title, &item.Description, &item.CreatedAt)
			if err != nil {
				http.Error(w, "Could not read invites", http.StatusInternalServerError)
				return
			}
			groups = append(groups, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}

func ListGroupJoinRequestsHandler(db *sql.DB) http.HandlerFunc {
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

		groupID := r.URL.Query().Get("group_id")
		if groupID == "" {
			http.Error(w, "Missing group_id", http.StatusBadRequest)
			return
		}

		var creatorID string
		err = db.QueryRow("SELECT creator_id FROM groups WHERE id = ?", groupID).Scan(&creatorID)
		if err != nil {
			http.Error(w, "Group not found", http.StatusNotFound)
			return
		}

		if creatorID != currentUserID {
			http.Error(w, "Only group creator can view join requests", http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT u.id, u.first_name, u.last_name, u.nickname
			FROM group_members gm
			JOIN users u ON gm.user_id = u.id
			WHERE gm.group_id = ? AND gm.status = 'requested'
			ORDER BY gm.created_at DESC
		`, groupID)
		if err != nil {
			http.Error(w, "Could not fetch join requests", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type GroupJoinRequestUser struct {
			UserID    string `json:"user_id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Nickname  string `json:"nickname"`
		}

		users := make([]GroupJoinRequestUser, 0)

		for rows.Next() {
			var item GroupJoinRequestUser
			var nickname sql.NullString

			err := rows.Scan(&item.UserID, &item.FirstName, &item.LastName, &nickname)
			if err != nil {
				http.Error(w, "Could not read join requests", http.StatusInternalServerError)
				return
			}

			if nickname.Valid {
				item.Nickname = nickname.String
			} else {
				item.Nickname = ""
			}

			users = append(users, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func ListMyGroupsHandler(db *sql.DB) http.HandlerFunc {
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
			SELECT g.id, g.creator_id, g.title, g.description, g.created_at,
			       (SELECT COUNT(*) FROM group_members WHERE group_id = g.id AND status = 'accepted') AS member_count,
			       'accepted' AS membership_status
			FROM group_members gm
			JOIN groups g ON gm.group_id = g.id
			WHERE gm.user_id = ? AND gm.status = 'accepted'
			ORDER BY g.created_at DESC
		`, currentUserID)
		if err != nil {
			http.Error(w, "Could not fetch user groups", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		groups := make([]GroupResponse, 0)

		for rows.Next() {
			var item GroupResponse
			err := rows.Scan(
				&item.ID,
				&item.CreatorID,
				&item.Title,
				&item.Description,
				&item.CreatedAt,
				&item.MemberCount,
				&item.MembershipStatus,
			)
			if err != nil {
				http.Error(w, "Could not read user groups", http.StatusInternalServerError)
				return
			}
			groups = append(groups, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(groups)
	}
}

func isAcceptedGroupMember(db *sql.DB, groupID, userID string) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM group_members
		WHERE group_id = ? AND user_id = ? AND status = 'accepted'
	`, groupID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func handleGroupMembershipStatusUpdate(db *sql.DB, w http.ResponseWriter, r *http.Request, currentStatus, newStatus, successMessage string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUserID, err := auth.GetCurrentUserID(db, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body GroupActionRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.GroupID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE group_members
		SET status = ?
		WHERE group_id = ? AND user_id = ? AND status = ?
	`, newStatus, body.GroupID, currentUserID, currentStatus)
	if err != nil {
		http.Error(w, "Could not update group membership", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Matching group membership record not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(successMessage))
}

func handleCreatorGroupRequestStatusUpdate(db *sql.DB, w http.ResponseWriter, r *http.Request, currentStatus, newStatus, successMessage string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentUserID, err := auth.GetCurrentUserID(db, r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var body GroupActionRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.GroupID == "" || body.UserID == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var creatorID string
	err = db.QueryRow("SELECT creator_id FROM groups WHERE id = ?", body.GroupID).Scan(&creatorID)
	if err != nil {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if creatorID != currentUserID {
		http.Error(w, "Only group creator can manage join requests", http.StatusForbidden)
		return
	}

	result, err := db.Exec(`
		UPDATE group_members
		SET status = ?
		WHERE group_id = ? AND user_id = ? AND status = ?
	`, newStatus, body.GroupID, body.UserID, currentStatus)
	if err != nil {
		http.Error(w, "Could not update join request", http.StatusInternalServerError)
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, "Matching join request not found", http.StatusNotFound)
		return
	}

	if newStatus == "accepted" {
		var groupTitle string
		_ = db.QueryRow("SELECT title FROM groups WHERE id = ?", body.GroupID).Scan(&groupTitle)
		_ = createNotification(
			db,
			body.UserID,
			"group_accept",
			body.GroupID,
			"Your request to join group "+groupTitle+" was accepted",
		)
	}

	w.Write([]byte(successMessage))
}
