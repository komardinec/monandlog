package main

import (
	"net/http"
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
