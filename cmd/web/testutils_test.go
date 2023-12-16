package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"tarala/snippetbox/pkg/models/mock"
	"testing"
	"time"

	"github.com/golangcollege/sessions"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {

	templateCache, err := newTemplateCache("./../../ui/html")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
		//thanks to fact that application struct has interfaces as type, we can pass more than one implementation here
		snippets: &mock.SnippetModel{},
		users: &mock.UserModel{}		,
		templateCache: templateCache,
		session: session,
		
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server
// instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	//cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	//assign new cookieJar to our test client
	//this is going to store all cookies server saves during request processing
	//we can test cookies that way
	ts.Client().Jar = jar


	//disable following redirect requests
	//this function is called, if 3xx response is received by a client, 
	// and it returns ErrUseLastResponse forces to return the receiver response, not follow redirect 
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Implement a get method on our custom testServer type. This makes a GET
// request to a given url path on the test server, and returns the response
// status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	//using test server client to send http request
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return rs.StatusCode, rs.Header, body
}
