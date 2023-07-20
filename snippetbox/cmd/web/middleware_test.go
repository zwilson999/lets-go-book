package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.lets-go/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a mock HTTP handler that we can pass to our secureHeaders middleware
	// which writes a 200 status code and "OK"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// pass the mock handler to our middleware
	// because secureHeaders returns a http.Handler we can call its ServeHTTP() method
	// by passing in the http.ResponseRecorder and dummy http.Request
	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// check the middleware has correctly set the Referrer-Policy header
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// check the middleware has the correct X-Content-Type-Options header
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// check the middleware has the correct X-Frame-Options header
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	// check the middle ware has the next handler called in line and the response statuses are as expected
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	// assert our response bodys are equivalent
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
