package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"social-network/pkg/auth"

	"github.com/google/uuid"
)

type CreateEventRequest struct {
	GroupID       string `json:"group_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	EventDateTime string `json:"event_datetime"`
}

type EventResponseRequest struct {
	EventID  string `json:"event_id"`
	Response string `json:"response"`
}

type GroupEventResponse struct {
	ID            string `json:"id"`
	GroupID       string `json:"group_id"`
	CreatorID     string `json:"creator_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	EventDateTime string `json:"event_datetime"`
	CreatedAt     string `json:"created_at"`
}

type EventMemberResponse struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Response  string `json:"response"`
	CreatedAt string `json:"created_at"`
}

func CreateGroupEventHandler(db *sql.DB) http.HandlerFunc {
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

		var body CreateEventRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.GroupID == "" || body.Title == "" || body.Description == "" || body.EventDateTime == "" {
			http.Error(w, "group_id, title, description and event_datetime are required", http.StatusBadRequest)
			return
		}
		body.Title = strings.TrimSpace(body.Title)
		body.Description = strings.TrimSpace(body.Description)
		if len(body.Title) > 100 {
			http.Error(w, "Event title cannot exceed 100 characters", http.StatusBadRequest)
			return
		}
		if len(body.Description) > 500 {
			http.Error(w, "Event description cannot exceed 500 characters", http.StatusBadRequest)
			return
		}
		eventDT, dtErr := time.Parse("2006-01-02T15:04", body.EventDateTime)
		if dtErr != nil {
			eventDT, dtErr = time.Parse("2006-01-02T15:04:05", body.EventDateTime)
			if dtErr != nil {
				http.Error(w, "Invalid event date/time format", http.StatusBadRequest)
				return
			}
		}
		if !eventDT.After(time.Now()) {
			http.Error(w, "Event date must be in the future", http.StatusBadRequest)
			return
		}

		ok, err := isAcceptedGroupMemberHandler(db, body.GroupID, currentUserID)
		if err != nil || !ok {
			http.Error(w, "Only group members can create events", http.StatusForbidden)
			return
		}

		eventID := uuid.New().String()

		_, err = db.Exec(`
			INSERT INTO events (id, group_id, creator_id, title, description, event_datetime)
			VALUES (?, ?, ?, ?, ?, ?)
		`, eventID, body.GroupID, currentUserID, body.Title, body.Description, body.EventDateTime)
		if err != nil {
			http.Error(w, "Could not create event", http.StatusInternalServerError)
			return
		}

		var creatorFirstName string
		var creatorLastName string
		_ = db.QueryRow(
			"SELECT first_name, last_name FROM users WHERE id = ?",
			currentUserID,
		).Scan(&creatorFirstName, &creatorLastName)

		rows, err := db.Query(`
			SELECT user_id
			FROM group_members
			WHERE group_id = ? AND status = 'accepted' AND user_id != ?
		`, body.GroupID, currentUserID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var memberID string
				if rows.Scan(&memberID) == nil {
					_ = createNotification(
						db,
						memberID,
						"group_event_created",
						eventID,
						creatorFirstName+" "+creatorLastName+" created a new event: "+body.Title,
					)
				}
			}
		}

		w.Write([]byte("Group event created successfully"))
	}
}

func ListGroupEventsHandler(db *sql.DB) http.HandlerFunc {
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
			http.Error(w, "group_id required", http.StatusBadRequest)
			return
		}

		ok, err := isAcceptedGroupMemberHandler(db, groupID, currentUserID)
		if err != nil || !ok {
			http.Error(w, "Not a group member", http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT id, group_id, creator_id, title, description, event_datetime, created_at
			FROM events
			WHERE group_id = ?
			ORDER BY event_datetime ASC
		`, groupID)
		if err != nil {
			http.Error(w, "Could not fetch events", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		events := make([]GroupEventResponse, 0)

		for rows.Next() {
			var item GroupEventResponse
			err := rows.Scan(
				&item.ID,
				&item.GroupID,
				&item.CreatorID,
				&item.Title,
				&item.Description,
				&item.EventDateTime,
				&item.CreatedAt,
			)
			if err != nil {
				http.Error(w, "Could not read events", http.StatusInternalServerError)
				return
			}
			events = append(events, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}
}

func RespondToEventHandler(db *sql.DB) http.HandlerFunc {
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

		var body EventResponseRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.EventID == "" || body.Response == "" {
			http.Error(w, "event_id and response are required", http.StatusBadRequest)
			return
		}

		if body.Response != "going" && body.Response != "not_going" {
			http.Error(w, "response must be 'going' or 'not_going'", http.StatusBadRequest)
			return
		}

		var groupID string
		err = db.QueryRow(`
			SELECT group_id
			FROM events
			WHERE id = ?
		`, body.EventID).Scan(&groupID)
		if err != nil {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		ok, err := isAcceptedGroupMemberHandler(db, groupID, currentUserID)
		if err != nil || !ok {
			http.Error(w, "Only group members can respond to events", http.StatusForbidden)
			return
		}

		var existingID string
		err = db.QueryRow(`
			SELECT id
			FROM event_responses
			WHERE event_id = ? AND user_id = ?
		`, body.EventID, currentUserID).Scan(&existingID)

		if err == sql.ErrNoRows {
			_, err = db.Exec(`
				INSERT INTO event_responses (id, event_id, user_id, response)
				VALUES (?, ?, ?, ?)
			`, uuid.New().String(), body.EventID, currentUserID, body.Response)
			if err != nil {
				http.Error(w, "Could not save event response", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("Event response saved"))
			return
		}

		if err != nil {
			http.Error(w, "Could not check existing response", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			UPDATE event_responses
			SET response = ?
			WHERE event_id = ? AND user_id = ?
		`, body.Response, body.EventID, currentUserID)
		if err != nil {
			http.Error(w, "Could not update event response", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Event response updated"))
	}
}

func ListEventResponsesHandler(db *sql.DB) http.HandlerFunc {
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

		eventID := r.URL.Query().Get("event_id")
		if eventID == "" {
			http.Error(w, "event_id required", http.StatusBadRequest)
			return
		}

		var groupID string
		err = db.QueryRow(`
			SELECT group_id
			FROM events
			WHERE id = ?
		`, eventID).Scan(&groupID)
		if err != nil {
			http.Error(w, "Event not found", http.StatusNotFound)
			return
		}

		ok, err := isAcceptedGroupMemberHandler(db, groupID, currentUserID)
		if err != nil || !ok {
			http.Error(w, "Not a group member", http.StatusForbidden)
			return
		}

		rows, err := db.Query(`
			SELECT er.user_id, u.first_name, u.last_name, er.response, er.created_at
			FROM event_responses er
			LEFT JOIN users u ON u.id = er.user_id
			WHERE er.event_id = ?
			ORDER BY er.created_at ASC
		`, eventID)
		if err != nil {
			http.Error(w, "Could not fetch event responses", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		responses := make([]EventMemberResponse, 0)

		for rows.Next() {
			var item EventMemberResponse
			err := rows.Scan(&item.UserID, &item.FirstName, &item.LastName, &item.Response, &item.CreatedAt)
			if err != nil {
				http.Error(w, "Could not read event responses", http.StatusInternalServerError)
				return
			}
			responses = append(responses, item)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
	}
}
