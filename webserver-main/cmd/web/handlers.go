package main

import (
	"fmt"
	"net/http"
)

// welcomeTest handler function which writes a byte slice containing a string
// as the response body
func WelcomeTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home, komardinec\n"))
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
