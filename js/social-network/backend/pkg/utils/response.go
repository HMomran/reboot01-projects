package utils

import (
	"encoding/json"
	"net/http"
)

// JSON sends a full JSON response
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Message sends a success message
func Message(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{
		"message": msg,
	})
}

// Error sends an error message
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{
		"error": msg,
	})
}
