package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func secureHeaders(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("X-XSS-Protection", "1; mode=block")
// 		w.Header().Set("X-Frame-Options", "deny")

// 		next.ServeHTTP(w, r)
// 	}

// 	return http.HandlerFunc(fn)
// }

func TestSecureHeaders(t *testing.T) {
	//enable parallel processing of this test along with the other tests
	//this can boost the time spent on testing
	t.Parallel()
	rr := httptest.NewRecorder()
	r, err := http.NewRequest("GET","/",nil)
	if err != nil {
		t.Fatal(err)
	}

	//create mock HTTP Handler which writes a 200 status and OK to response body
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	//check headers
	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}


	//check body
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body equal to %q", "OK")
	}
}