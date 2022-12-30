package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// This serverError helper writes an error message and stack trace to errorLog
// then sends a generic 500 Internal Server Error response to user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) // debug.Stack() gets stack trace for current goroutines
	app.errorLog.Output(2, trace)                              // Change depth of stack trace

	// Write internal server error to response writer
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// This helper sends specific status codes to the user as well as the status text
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Also construct a clientError wrapper for statusNotFound which will send a 404
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
