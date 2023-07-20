package main

import (
	"net/http"
	"testing"

	"snippetbox.lets-go/internal/assert"
)

func TestPing(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	code, _, respBody := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, respBody, "OK")
}
