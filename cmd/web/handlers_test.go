package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// create recorder which replaces writer
	rr := httptest.NewRecorder()

	// prepare http request struct

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// invoke tested method
	ping(rr, r)

	// check the results
	rs := rr.Result()

	// check if status is ok
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// to check request body we need more code
	// as body is a byte array
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal()
	}

	if string(body) != "OK" {
		t.Errorf("want %v", "OK")
	}
}

func TestPingEToE(t *testing.T) {
	//I would call it "test context"
	app := newTestApplication(t)

	//This is test TLS server. We use routes to run exact functionality on test server
	// The server starts on random URL and PORT. You can access them in ts.URL variable
	ts:= newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want %v", "OK")
	}
}
