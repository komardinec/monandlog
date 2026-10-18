package main

import (
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

	// Create dynamicMiddlewareChain specific to our
	// dynamic application routes.
	dynamicMiddlewareChain := alice.New(app.sessionManager.LoadAndSave)

	// Wrap handle in dynamicMiddlewareChain, notice that because .ThenFunc()
	// returns a http.Handler (rather then http.HandlerFunc) we need to set mux.Handle (instead of mux.HandleFunc)
	mux.Handle("GET /{$}", dynamicMiddlewareChain.ThenFunc(app.HomePage)) // restrict catch-all to match "/" only
	mux.Handle("GET /file/view/{id}", dynamicMiddlewareChain.ThenFunc(app.ViewFile))
	mux.Handle("GET /file/view-latest", dynamicMiddlewareChain.ThenFunc(app.Latest))
	mux.Handle("GET /file/create", dynamicMiddlewareChain.ThenFunc(app.CreateFile))
	mux.Handle("POST /file/create", dynamicMiddlewareChain.ThenFunc(app.CreateFilePost))

	// Create middlewareChain which to be applied for every request application receives
	staticMiddlewareChain := alice.New(app.RecoverPanic, app.ConnectionLogger, CommonHeaders)

	// Apply middlewareChain to servemux
	return staticMiddlewareChain.Then(mux)
}
