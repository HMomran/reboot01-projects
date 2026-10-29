package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"real-time-forum/db"
	"real-time-forum/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Nickname  string `json:"nickname"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Nickname = strings.TrimSpace(req.Nickname)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Gender = strings.TrimSpace(req.Gender)

	if req.Nickname == "" || req.FirstName == "" || req.LastName == "" ||
		req.Email == "" || req.Gender == "" || req.Password == "" {
		jsonError(w, "all fields are required", http.StatusBadRequest)
		return
	}

	// --- field-level validation ---
	if msg := validateNickname(req.Nickname); msg != "" {
		jsonError(w, msg, http.StatusBadRequest)
		return
	}
	if msg := validateName("first name", req.FirstName); msg != "" {
		jsonError(w, msg, http.StatusBadRequest)
		return
	}
	if msg := validateName("last name", req.LastName); msg != "" {
		jsonError(w, msg, http.StatusBadRequest)
		return
	}
	if msg := validateEmail(req.Email); msg != "" {
		jsonError(w, msg, http.StatusBadRequest)
		return
	}
	if req.Age < minAge || req.Age > maxAge {
		jsonError(w, "age must be between 13 and 120", http.StatusBadRequest)
		return
	}
	if msg := validateGender(req.Gender); msg != "" {
		jsonError(w, msg, http.StatusBadRequest)
		return
	}
	if msg := validatePassword(req.Password); msg != "" {
		jsonError(w, msg, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:        uuid.NewString(),
		Nickname:  req.Nickname,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Age:       req.Age,
		Gender:    req.Gender,
		Password:  string(hash),
	}

	if err := db.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "nickname") {
				jsonError(w, "nickname already taken", http.StatusConflict)
			} else {
				jsonError(w, "email already registered", http.StatusConflict)
			}
			return
		}
		jsonError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	jsonOK(w, http.StatusCreated, map[string]string{"message": "registration successful"})
}
