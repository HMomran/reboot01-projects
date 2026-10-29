package backend

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
)

func RenderErrorPage(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	tmpl, err := template.ParseFiles("Frontend/templates/error/ErrorPage.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, map[string]interface{}{
		"Code":    code,
		"Message": message,
	})
}

func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format (must be abc@abc.abc)")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	upper := regexp.MustCompile(`[A-Z]`)
	lower := regexp.MustCompile(`[a-z]`)
	number := regexp.MustCompile(`[0-9]`)

	if !upper.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !lower.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !number.MatchString(password) {
		return errors.New("password must contain at least one digit")
	}
	return nil
}