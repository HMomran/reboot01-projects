package auth

import (
    "forum/backend"
    "forum/backend/database"
    "html/template"
    "net/http"
    "time"
    "golang.org/x/crypto/bcrypt"
    "net/url"
)

func SignIn(w http.ResponseWriter, r *http.Request, DB *database.DB) {
    cookie, err := r.Cookie("session_token")
    if err == nil {
        _, err := DB.GetUserIDFromSession(cookie.Value)
        if err == nil {
            http.Redirect(w, r, "/communication", http.StatusSeeOther)
            return
        }
    }

    if r.Method == http.MethodPost {
        email := r.FormValue("email")
        password := r.FormValue("password")

        id, hashedPassword, err := DB.GetUserByEmail(email)
        if err != nil {
            errorMsg := url.QueryEscape("Invalid email or password")
            http.Redirect(w, r, "/signin?error="+errorMsg, http.StatusSeeOther)
            return
        }

        err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
        if err != nil {
            errorMsg := url.QueryEscape("Invalid email or password")
            http.Redirect(w, r, "/signin?error="+errorMsg, http.StatusSeeOther)
            return
        }

        err = DB.DeleteUserSessions(id)
        if err != nil {
            backend.RenderErrorPage(w, http.StatusInternalServerError, "Error removing old session")
            return
        }

        sessionID, expiry, err := DB.CreateSession(id, 2*time.Hour)
        if err != nil {
            backend.RenderErrorPage(w, http.StatusInternalServerError, "Error creating session")
            return
        }

        http.SetCookie(w, &http.Cookie{
            Name:    "session_token",
            Value:   sessionID,
            Expires: expiry,
            Path:    "/",
        })

        http.Redirect(w, r, "/communication", http.StatusSeeOther)
        return
    }

    errorMsg := r.URL.Query().Get("error")
    data := struct {
        Error string
    }{
        Error: errorMsg,
    }

    signin, err := template.ParseFiles("Frontend/templates/SignIn.html")
    if err != nil {
        backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading page")
        return
    }
    signin.Execute(w, data)
}
