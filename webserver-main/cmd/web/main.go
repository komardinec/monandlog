package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Define destinct parameter variable to control the server port number
	serverPort := 8080

	// Use http.NewServeMux function to initialize a new servemux(controller),
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory
	// Path are relative to webserver directory root (execution place)
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register a file server as the handler
	// for all URL paths starting with the path given in http.Dir() function
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer)) // http.StripPrefix() will remove "/static" from a filesystem path
	// to make possible looking for a files inside .ui/static/ sub-directory

	// then register WelcomeTest as a handler for root pattern "/"
	mux.HandleFunc("GET /{$}", HomePage) // restrict catch-all to match "/" only
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
