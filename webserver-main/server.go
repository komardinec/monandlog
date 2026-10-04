package main

import (
	"fmt"
	"log"
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

func main() {
	// Define destinct parameter variable to control the server port number
	serverPort := 8080

	// Use http.NewServeMux function to initialize a new servemux(controller),
	mux := http.NewServeMux()
	// then register WelcomeTest as a handler for root pattern "/"
	mux.HandleFunc("GET /{$}", WelcomeTest) // restrict catch-all to match / only
	mux.HandleFunc("GET /file/view/{filename}", ViewFile)
	mux.HandleFunc("GET /file/create", CreateFile)
	mux.HandleFunc("POST /file/create", CreateFilePost)

	// Print starting log message
	log.Printf("Server is running on port %d. Ctrl+C to stop it...", serverPort)

	// Use the http.ListenAndServe() function starts a new web server.
	// We pass inside 2 arguments - TCP address to listen on
	// and servemux created earlier.
	// Based on information from documentation http.ListenAndServe always return
	// non-nil error so if we something goes wrong log.Fatal()
	// will print it as a fatal error
	err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), mux)
	log.Fatal(err)
}
