package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper which makes an instance of our app struct for mocked dependencies
func newTestApplication(t *testing.T) *application {
	return &application{
		infoLog:  log.New(io.Discard, "", 0),
		errorLog: log.New(io.Discard, "", 0),
	}
}

// embed a httptest.Server within this testServer struct
type testServer struct {
	*httptest.Server
}

// helper to create a new test server which returns one of our custom testServer structs
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

// implement a get() method on our custom testServer type
// this will make a GET request to a given url path using the test server client, and returns
// the response status code, headers and body
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	resp, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(respBody)
	return resp.StatusCode, resp.Header, string(respBody)

}
