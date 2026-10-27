package main

import (
	"net/http"

	"debutwebapp.komardinec/ui"

	"github.com/justinas/alice"
)

// Update the signature for the Route() method, so that
// now returns http.Handler instead of *http.ServeMux
func (app *application) Route() http.Handler {
	// Use http.NewServeMux function to initialize a new servemux(controller),
	mux := http.NewServeMux()

	// Use the http.FileServerFS() function to create a HTTP handler which
	// serves the embedded files in ui.Files. It's important to note that our
	// static files are contained in the "static" folder of the ui.Files
	// embedded filesystem. So, for examples, our CSS stylesheet is located at
	// "static/css/main.css". This means that we no longer need to strip the
	// prefix from the request URL - any request that start with /static/ can
	// just be passed directly to the file server and the corresponding static
	// file will be served (so long as it exists).
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

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
