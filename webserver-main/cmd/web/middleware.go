package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func CommonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Go-Version", "1.27.x")

		next.ServeHTTP(w, r)
	})
}

func (app *application) ConnectionLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			req_ip     = r.RemoteAddr
			req_proto  = r.Proto
			req_method = r.Method
			req_uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", req_ip, "proto", req_proto, "method", req_method, "uri", req_uri)

		next.ServeHTTP(w, r)
	})
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

func (app *application) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page
		// return from the middleware chain so that no subsequent handlers in
		// the chain are executed.
		if !app.IsAuthenticated(r) {
			app.sessionManager.Put(r.Context(), "flash", "We don't know who you are. Enter into your account first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no store" header so that pages
		// require authentication are not stored in the users browser cache
		// (or other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// And call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}

func (app *application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authentificatedUserID value from the session using the
		// GetInt() method. This will return the zero value for an int (0) if no
		// "authentificatedUserID" value in the session -- in which case we
		// call the next handler in the chain as normal and return
		id := app.sessionManager.GetInt(r.Context(), "authentificatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise, we check to see if a user with that ID exists in
		// our database
		exists, err := app.users.Exists(id)
		if err != nil {
			app.ServerError(w, r, err)
			return
		}

		// If a matching user is found, we know that the request is
		// coming from an authenticated user who exists if our database
		// create a new copy of the request (with an IsAuthentificatedContextKey
		// value of true in the request context) and assign it to r
		if exists {
			ctx := context.WithValue(r.Context(), IsAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
