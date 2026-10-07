package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// It is not recommended to use global variables at general,
// but if we want to use templates for all handlers
// global []string are very good implementation
var templateFiles = []string{
	"./ui/html/base.tmpl.html",
	"./ui/html/pages/home.tmpl.html",
	"./ui/html/partials/nav.tmpl.html",
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "GO v1.27.1")

	// template.ParseFiles will read the template file and turn its contents
	// into a template set. If there is an error occured we will get into terminal the error message
	// and user will get nice little plain text "Internal Server Error" response
	// and that case handler will return and no subsequent code will be executed.
	ts, err := template.ParseFiles(templateFiles...) // we unpacking string slice of templateFiles paths
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// ExecuteTemplate() method will write contents of the "base"
	// template as the response body
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Handler for viewing the files (test function, meant to be depricated later)
func ViewFile(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")

	if filename == "" {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Currently looking at: %s", filename)
}

// Handler for creating the files (test function, meant to be depricated later)
func CreateFile(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("File creation form"))
}

// Handler to work with POST method while creating a new file (test function, meant to be depricated later)
func CreateFilePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Test Go Backend Server")
	w.Header().Add("Author", "komardinec")

	// Send 201 status code
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Creating new file..."))
}
