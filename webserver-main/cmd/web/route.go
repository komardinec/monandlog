package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
)

// Update the signature for the Route() method, so that
// now returns http.Handler instead of *http.ServeMux
func (app *application) Route() http.Handler {
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
	mux.HandleFunc("GET /file/view/{id}", app.ViewFile)
	mux.HandleFunc("GET /file/view/latest", app.Latest)
	mux.HandleFunc("GET /file/create", app.CreateFile)
	mux.HandleFunc("POST /file/create", app.CreateFilePost)

	// Create middlewareChain which to be applied for every request application receives
	middlewareChain := alice.New(app.RecoverPanic, app.ConnectionLogger, CommonHeaders)

	// Apply middlewareChain to servemux
	return middlewareChain.Then(mux)
}

func (app *application) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Differed function will execute if pamic will occur in code,
		// and will always run as Go unwinds the call stack
		defer func() {
			// Built-in recover function will check if there has been a panic
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call of the ServerError method will return a 500 Internal Server Error
				app.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
