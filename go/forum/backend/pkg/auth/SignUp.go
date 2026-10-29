package auth

import (
	"forum/backend"
	"forum/backend/database"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

func SignUp(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		_, err := DB.GetUserIDFromSession(cookie.Value)
		if err == nil {
			http.Redirect(w, r, "/communication", http.StatusSeeOther)
			return
		}
	}

	if r.Method == http.MethodPost {
		name := strings.TrimSpace(r.FormValue("userName"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := strings.TrimSpace(r.FormValue("password"))

		if name == "" || email == "" || password == "" {
			redirectURL := "/signup?" +
				"error=" + url.QueryEscape("All fields are required") +
				"&userName=" + url.QueryEscape(name) +
				"&email=" + url.QueryEscape(email)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		if err := backend.ValidateEmail(email); err != nil {
			redirectURL := "/signup?" +
				"error=" + url.QueryEscape(err.Error()) +
				"&userName=" + url.QueryEscape(name) +
				"&email=" + url.QueryEscape(email)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		if err := backend.ValidatePassword(password); err != nil {
			redirectURL := "/signup?" +
				"error=" + url.QueryEscape(err.Error()) +
				"&userName=" + url.QueryEscape(name) +
				"&email=" + url.QueryEscape(email)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		exists, err := DB.CheckEmailExists(email)
		if err != nil {
			redirectURL := "/signup?" +
				"error=" + url.QueryEscape("Database error") +
				"&userName=" + url.QueryEscape(name) +
				"&email=" + url.QueryEscape(email)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		if exists {
			redirectURL := "/signup?" +
				"error=" + url.QueryEscape("Email already registered") +
				"&userName=" + url.QueryEscape(name) +
				"&email=" + url.QueryEscape(email)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		err = DB.CreateUser(name, email, password)
		if err != nil {
			redirectURL := "/signup?" +
				"error=" + url.QueryEscape("Failed to register user") +
				"&userName=" + url.QueryEscape(name) +
				"&email=" + url.QueryEscape(email)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	errorMsg := r.URL.Query().Get("error")
	email := r.URL.Query().Get("email")
	userName := r.URL.Query().Get("userName")

	data := struct {
		Error    string
		Email    string
		UserName string
	}{
		Error:    errorMsg,
		Email:    email,
		UserName: userName,
	}

	signup, err := template.ParseFiles("Frontend/templates/SignUp.html")
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading page")
		return
	}
	signup.Execute(w, data)
}
