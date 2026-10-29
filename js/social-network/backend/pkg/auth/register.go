package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	Avatar      string `json:"avatar"`
	Nickname    string `json:"nickname"`
	AboutMe     string `json:"about_me"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		req.Email = strings.TrimSpace(req.Email)
		req.FirstName = strings.TrimSpace(req.FirstName)
		req.LastName = strings.TrimSpace(req.LastName)
		req.Nickname = strings.TrimSpace(req.Nickname)
		req.AboutMe = strings.TrimSpace(req.AboutMe)

		if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.DateOfBirth == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		emailRe := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
		if !emailRe.MatchString(req.Email) {
			http.Error(w, "Invalid email address", http.StatusBadRequest)
			return
		}

		if len(req.FirstName) < 2 || len(req.FirstName) > 50 {
			http.Error(w, "First name must be between 2 and 50 characters", http.StatusBadRequest)
			return
		}
		if len(req.LastName) < 2 || len(req.LastName) > 50 {
			http.Error(w, "Last name must be between 2 and 50 characters", http.StatusBadRequest)
			return
		}

		if len(req.Password) < 8 {
			http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
			return
		}
		if len(req.Password) > 72 {
			http.Error(w, "Password cannot exceed 72 characters", http.StatusBadRequest)
			return
		}

		dob, dobErr := time.Parse("2006-01-02", req.DateOfBirth)
		if dobErr != nil {
			http.Error(w, "Invalid date of birth (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		minBirthDate := time.Now().AddDate(-13, 0, 0)
		if dob.After(minBirthDate) {
			http.Error(w, "You must be at least 13 years old to register", http.StatusBadRequest)
			return
		}

		if req.Nickname == "" {
			req.Nickname = generateUniqueNickname(db, req.FirstName, req.LastName)
		} else {
			if len(req.Nickname) > 30 {
				http.Error(w, "Nickname cannot exceed 30 characters", http.StatusBadRequest)
				return
			}
			nickFormatRe := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
			if !nickFormatRe.MatchString(req.Nickname) {
				http.Error(w, "Nickname can only contain letters, numbers, and underscores", http.StatusBadRequest)
				return
			}
			taken, err := isNicknameTaken(db, req.Nickname, "")
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			if taken {
				http.Error(w, "Nickname is already taken. Please choose another one.", http.StatusBadRequest)
				return
			}
		}

		if len(req.AboutMe) > 500 {
			http.Error(w, "About me cannot exceed 500 characters", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Could not hash password", http.StatusInternalServerError)
			return
		}

		userID := uuid.New().String()

		_, err = db.Exec(`
			INSERT INTO users (
				id, email, password, first_name, last_name, date_of_birth, avatar, nickname, about_me
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			userID,
			req.Email,
			string(hashedPassword),
			req.FirstName,
			req.LastName,
			req.DateOfBirth,
			nullIfEmpty(req.Avatar),
			nullIfEmpty(req.Nickname),
			nullIfEmpty(req.AboutMe),
		)
		if err != nil {
			http.Error(w, "Account already exists with this email. Please log in instead.", http.StatusBadRequest)
			return
		}

		err = CreateSession(db, userID, w)
		if err != nil {
			http.Error(w, "Could not create session", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User registered successfully"))
	}
}

func nullIfEmpty(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func isNicknameTaken(db *sql.DB, nickname string, excludeUserID string) (bool, error) {
	var count int
	var err error
	if excludeUserID != "" {
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(nickname) = LOWER(?) AND id != ?", nickname, excludeUserID).Scan(&count)
	} else {
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(nickname) = LOWER(?)", nickname).Scan(&count)
	}
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func generateUniqueNickname(db *sql.DB, firstName, lastName string) string {
	base := strings.ToLower(firstName + "_" + lastName)
	reg := regexp.MustCompile("[^a-z0-9_]")
	base = reg.ReplaceAllString(base, "")
	if base == "" || base == "_" {
		base = "user"
	}

	nickname := base
	suffix := 1
	for {
		taken, err := isNicknameTaken(db, nickname, "")
		if err == nil && !taken {
			return nickname
		}
		if len(base) > 20 {
			nickname = fmt.Sprintf("%s_%d", base[:20], suffix)
		} else {
			nickname = fmt.Sprintf("%s_%d", base, suffix)
		}
		suffix++
	}
}
