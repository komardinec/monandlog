package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"debutwebapp.komardinec/internal/models/mocks"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/justinas/alice"
)

// NewTestApplication() will return an instance of
// an application struct with mocked dependencies
func NewTestApplication(t *testing.T) *application {
	// Create an instance if the template cache
	templateCache, err := NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	// Add form decoder
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 1 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.DiscardHandler),
		dbWebApp:       &mocks.WebAppModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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

	// Disable redirect-following for the test server client by setting a custom
	// CheckRedirect function. This function will be called whenever a 3xx response
	// is received by the client, and by always returning a http.ErrUseLastResponse
	// error it forces the client to immediately return the receiver response.
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

// As an alternative approach, to not disturb the already compiled production router,
// we can include specific test routes to test the logic, because server logic available
// no matter what router we use
func (app *application) TestRoute() http.Handler {
	testMux := http.NewServeMux()

	generalMiddlewareChain := alice.New(app.RecoverPanic, app.ConnectionLogger, CommonHeaders)
	dynamicMiddlewareChain := alice.New(app.sessionManager.LoadAndSave, NoSurf, app.Authenticate)

	// Every path in our test router will mimic real paths, this will reduce errors in
	// code understanding
	testMux.Handle("GET /file/view/{id}", dynamicMiddlewareChain.ThenFunc(app.ViewFile))
	testMux.Handle("GET /user/signup", dynamicMiddlewareChain.ThenFunc(app.UserSignup))
	testMux.Handle("POST /user/signup", dynamicMiddlewareChain.ThenFunc(app.UserSignupPost))

	return generalMiddlewareChain.Then(testMux)
}

var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)">`)

func ExtractCSRFToken(t *testing.T, body string) string {
	// Use the FindStringSubmatch method to extract the token from the HTML body.
	// Note that this returns an array with the entire matched pattern in the
	// first position, and the values if any captured data in the subsequent
	// positions
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token was found in the body")
	}

	return html.UnescapeString(matches[1])
}

// Create a PostForm method for sending POST requests to the test server. The final parameter to this
// method us a url.Values object which can contain any form data that we want to send in the req body
func (ts *testServer) PostForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	// Read the response body from the test server
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	// Return the response status, headers and body
	return rs.StatusCode, rs.Header, string(body)
}
