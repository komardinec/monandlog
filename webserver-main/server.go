package main

import (
	"fmt"
	"log"
	"net/http"
)

// welcomeTest handler function which writes a byte slice containing a string
// as the response body
func WelcomeTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome Home, komardinec"))
}

func main() {
	// Define destinct parameter variable to control the server port number
	serverPort := 8080

	// Use http.NewServeMux function to initialize a new servemux(controller),
	mux := http.NewServeMux()
	// then register WelcomeTest as a handler for root pattern "/"
	mux.HandleFunc("/", WelcomeTest)

	// Print starting log message
	log.Printf("Server is running on port %d. Ctrl+C to stop it...", serverPort)

	// Use the http.ListenAndServe() function starts a new web server.
	// We pass inside 2 arguments - TCP address to listen on
	// and servemux created earlier.
	// Based on information from documentation http.ListenAndServe always return
	err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), mux)
	log.Fatal(err)
}
