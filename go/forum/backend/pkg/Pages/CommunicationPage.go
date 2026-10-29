package pages

import (
	"fmt"
	"forum/backend"
	"forum/backend/database"
	"html/template"
	"log"
	"net/http"
)

func CommunicationPage(w http.ResponseWriter, r *http.Request, DB *database.DB) {
	r.ParseForm()
	categoryFilters := r.Form["category"]
	filterType := r.URL.Query().Get("filter")

	var posts []database.Post
	var err error
	var userID int
	var isLoggedIn bool

	cookie, err := r.Cookie("session_token")
	if err == nil {
		userID, err = DB.GetUserIDFromSession(cookie.Value)
		isLoggedIn = (err == nil)
	}

	if filterType == "created" && isLoggedIn {
		posts, err = DB.GetPostsByUser(userID)
	} else if filterType == "liked" && isLoggedIn {
		posts, err = DB.GetLikedPostsByUser(userID)
	} else if len(categoryFilters) > 0 {
		var categoryIDs []int
		for _, catStr := range categoryFilters {
			var catID int
			fmt.Sscanf(catStr, "%d", &catID)
			categoryIDs = append(categoryIDs, catID)
		}
		posts, err = DB.GetPostsByMultipleCategories(categoryIDs)
	} else {
		posts, err = DB.GetAllPosts()
	}

	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading posts")
		return
	}

	categories, err := DB.GetAllCategories()
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading categories")
		return
	}

	for i := range posts {
		comments, err := DB.GetCommentsByPost(posts[i].ID)
		if err != nil {
			log.Printf("Error loading comments for post %d: %v", posts[i].ID, err)
			continue
		}
		posts[i].Comments = comments
	}
	errorMsg := r.URL.Query().Get("error")
	data := struct {
		Posts               []database.Post
		Categories          []database.Category
		IsLoggedIn          bool
		UserID              int
		SelectedCategories  []string
		ErrorMsg            string
		CurrentFilter       string
	}{
		Posts:              posts,
		Categories:         categories,
		IsLoggedIn:         isLoggedIn,
		UserID:             userID,
		SelectedCategories: categoryFilters,
		ErrorMsg:           errorMsg,
		CurrentFilter:      filterType,
	}

	tmpl, err := template.ParseFiles("Frontend/templates/CommunicationPage.html")
	if err != nil {
		backend.RenderErrorPage(w, http.StatusInternalServerError, "Error loading page")
		return
	}
	tmpl.Execute(w, data)
}