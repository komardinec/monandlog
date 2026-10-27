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
	dynamicMiddlewareChain := alice.New(app.sessionManager.LoadAndSave, NoSurf, app.Authenticate)

	// Wrap handle in dynamicMiddlewareChain, notice that because .ThenFunc()
	// returns a http.Handler (rather then http.HandlerFunc) we need to set mux.Handle (instead of mux.HandleFunc)
	mux.Handle("GET /{$}", dynamicMiddlewareChain.ThenFunc(app.HomePage)) // restrict catch-all to match "/" only
	mux.Handle("GET /user/signup", dynamicMiddlewareChain.ThenFunc(app.UserSignup))
	mux.Handle("POST /user/signup", dynamicMiddlewareChain.ThenFunc(app.UserSignupPost))
	mux.Handle("GET /user/login", dynamicMiddlewareChain.ThenFunc(app.UserLogin))
	mux.Handle("POST /user/login", dynamicMiddlewareChain.ThenFunc(app.UserLoginPost))

	// Authentification only aplication routes will use this chain,
	// it contains everything inside dynamic chain + defensive authentification middleware
	protectedMiddlewareChain := dynamicMiddlewareChain.Append(app.RequireAuthentication)

	mux.Handle("GET /file/view/{id}", protectedMiddlewareChain.ThenFunc(app.ViewFile))
	mux.Handle("GET /file/view-latest", protectedMiddlewareChain.ThenFunc(app.Latest))
	mux.Handle("GET /file/create", protectedMiddlewareChain.ThenFunc(app.CreateFile))
	mux.Handle("POST /file/create", protectedMiddlewareChain.ThenFunc(app.CreateFilePost))
	mux.Handle("POST /user/logout", protectedMiddlewareChain.ThenFunc(app.UserLogoutPost))

	// Create generalMiddlewareChain which to be applied for every request application receives
	generalMiddlewareChain := alice.New(app.RecoverPanic, app.ConnectionLogger, CommonHeaders)

	// Apply middlewareChain to servemux
	return generalMiddlewareChain.Then(mux)
}
