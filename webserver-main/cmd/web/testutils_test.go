package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

// NewTestApplication() will return an instance of
// an application struct with mocked dependencies
func NewTestApplication(t *testing.T) *application {
	return &application{
		logger: slog.New(slog.DiscardHandler),
	}
}

// Custom testServer struct type with embedded httptest.Server instance
type testServer struct {
	*httptest.Server
}

// NewTestServer() return reference to intiialized
// testServer structure
func NewTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// Initialize the new cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server client. Any response
	// cookies will now be stored and sent with subsequent requests
	// when using this client
	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Get() test method will run ts.Client().Get to specified path on server
// returning response status code, headers and body
func (ts *testServer) Get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
