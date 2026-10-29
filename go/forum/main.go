package main

import (
	"forum/backend/database"
	pages "forum/backend/pkg/Pages"
	"forum/backend/pkg/auth"
	"forum/backend/pkg/handlers"
	"log"
	"net/http"
)

var DB *database.DB

func main() {
	var err error
	DB, err = database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer DB.Close()

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./Frontend/static/"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {pages.WelcomePage(w, r, DB)})
	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {auth.SignUp(w, r, DB)})
	mux.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {auth.SignIn(w, r, DB)})
	mux.HandleFunc("/communication", func(w http.ResponseWriter, r *http.Request) {pages.CommunicationPage(w, r, DB)})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {auth.Logout(w, r, DB)})
	mux.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {handlers.CreatePost(w, r, DB)})
	mux.HandleFunc("/create-post-page", func(w http.ResponseWriter, r *http.Request) {pages.CreatePostPage(w, r, DB)})
	mux.HandleFunc("/like", func(w http.ResponseWriter, r *http.Request) {handlers.LikePost(w, r, DB)})
	mux.HandleFunc("/dislike", func(w http.ResponseWriter, r *http.Request) {handlers.DislikePost(w, r, DB)})
	mux.HandleFunc("/add-comment", func(w http.ResponseWriter, r *http.Request) {handlers.AddComment(w, r, DB)})
	mux.HandleFunc("/like-comment", func(w http.ResponseWriter, r *http.Request) {handlers.LikeComment(w, r, DB)})
	mux.HandleFunc("/dislike-comment",func(w http.ResponseWriter, r *http.Request) {handlers.DislikeComment(w, r, DB)})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
// hello world