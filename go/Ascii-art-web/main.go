package main

import (
	"asciiart/functions"
	"fmt"
	"html/template"
	"net/http"
)

type Data struct {
	Result string
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	tmpl, err := template.ParseFiles("templates/error/400.html")
	if err != nil {
		http.Error(w, "400 - Bad Request: Input must contain only printable ASCII characters (32–126) and be no longer than 50", http.StatusBadRequest)
		return
	}
	tmpl.Execute(w, nil)
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)

	tmpl, err := template.ParseFiles("templates/error/500.html")
	if err != nil {
		http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {

	temp1, err := template.ParseFiles("templates/error/404.html")
	if err != nil {
		http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		return
	}

	temp1.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	temp1, err := template.ParseFiles("templates/index.html")
	if err != nil {
		internalServerErrorHandler(w, r)
		return
	}
	temp1.Execute(w, nil)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	text := r.FormValue("Data")
	banner := r.FormValue("banner")

	result := functions.AsciiArt(text, banner)

	if text == "" || banner == "" || result == "" {
		badRequestHandler(w, r)
		return
	}

	data := Data{
		Result: result,
	}

	temp1, err := template.ParseFiles("templates/result.html")
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	err = temp1.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		internalServerErrorHandler(w, r)
		return
	}

}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			homeHandler(w, r)
		} else if r.URL.Path == "/result" {
			resultHandler(w, r)
		} else {
			notFoundHandler(w, r)
		}
	})
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
