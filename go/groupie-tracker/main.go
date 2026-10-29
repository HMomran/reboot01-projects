package main

import (
	"fmt"
	"groupie/function"
	"html/template"
	"net/http"
	"strconv"
)

func ArtistListPage(w http.ResponseWriter, r *http.Request) {
	artists, err := function.FetchArtists()
	if err != nil {
		internalServerErrorHandler(w, r)
		return
	}

	tmpl, err := template.ParseFiles("backend/indix.html")
	if err != nil {
		fmt.Println(err)
		internalServerErrorHandler(w, r)
		return
	}
	tmpl.Execute(w, artists)
}

func ArtistDetailsPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		badRequestHandler(w, r)
		return
	}

	artists, err := function.FetchArtists()
	if err != nil {
		badRequestHandler(w, r)
		return
	}

	var selectedArtist *function.Artist
	for i := 0; i < len(artists); i++ {
		if artists[i].ID == id {
			selectedArtist = &artists[i]
			break
		}
	}

	if selectedArtist == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	relations, err := function.FetchRelation()
	if err != nil {
		badRequestHandler(w, r)
		return
	}

	var selectedRelation *function.Relation
	for i := 0; i < len(relations); i++ {
		if relations[i].ID == id {
			selectedRelation = &relations[i]
			break
		}
	}

	if selectedRelation == nil {
		selectedRelation = &function.Relation{ID: id, Dates: map[string][]string{}}
	}

	data := struct {
		Artist   *function.Artist
		Relation *function.Relation
	}{
		Artist:   selectedArtist,
		Relation: selectedRelation,
	}

	tmpl, err := template.ParseFiles("backend/details.html")
	if err != nil {
		notFoundHandler(w, r)
		return
	}
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			ArtistListPage(w, r)
		} else if r.URL.Path == "/details" {
			ArtistDetailsPage(w, r)
		} else {
			notFoundHandler(w, r)
		}
	})

	port := ":8080"
	fmt.Println("Server is running on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	temp1, err := template.ParseFiles("backend/error/404.html")
	if err != nil {
		internalServerErrorHandler(w, r)
		return
	}
	temp1.Execute(w, nil)
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl, err := template.ParseFiles("backend/error/400.html")
	if err != nil {
		internalServerErrorHandler(w, r)
		return
	}
	tmpl.Execute(w, nil)
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl, err := template.ParseFiles("backend/error/500.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
