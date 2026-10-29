package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"debutwebapp.komardinec/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	// Initialize a new httptest.ResponseRecorder and dummy http.Request
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our middleware
	// It will write the 200 code and some text
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am fish. Feed me some bread! Please?"))
	})

	// Pass the mock HTTP handler to our CommonHeaders middleware.
	// Because CommonHeaders returns a http.Handler we can call its ServeHTTP()
	// method, passing in the http.ResponseRecorder and dummy http.Request to
	// execute it
	CommonHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method on the http.ResponseRecorder to get the results
	// of the test
	rs := rr.Result()

	// Check the Content-Security-Policy header on the response
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// Check the Referrer-Policy header on the response
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// Check the X-Content-Type-Options header on the response
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// Check the X-Frame-Options header on the response
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	// Check the X-XSS-Protection header on the response
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	expectedValue = "1.27.x"
	assert.Equal(t, rs.Header.Get("Go-Version"), expectedValue)

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "I am fish. Feed me some bread! Please?")
}
