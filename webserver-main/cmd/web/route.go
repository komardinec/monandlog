package main

import "net/http"

// Separate functional method for servemux router, which is returning a reference to router in method
func (app *application) Route() *http.ServeMux {
	// Use http.NewServeMux function to initialize a new servemux(controller),
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory
	// Path are relative to webserver directory root (execution place)
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register a file server as the handler
	// for all URL paths starting with the path given in http.Dir() function
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer)) // http.StripPrefix() will remove "/static" from a filesystem path
	// to make possible looking for a files inside ./ui/static/ sub-directory

	// then register WelcomeTest as a handler for root pattern "/"
	mux.HandleFunc("GET /{$}", app.HomePage) // restrict catch-all to match "/" only
	mux.HandleFunc("GET /file/view/{filename}", app.ViewFile)
	mux.HandleFunc("GET /file/create", app.CreateFile)
	mux.HandleFunc("POST /file/create", app.CreateFilePost)

	return mux
}
