package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox!" as response body
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't then we will send a 404
	if r.URL.Path != "/" {
		app.notFound(w) // Helper for 404s
		return          // Important or else page will keep running
	}

	// Slice to contain our template files
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
		"./ui/html/partials/nav.tmpl",
	}

	// Use the template.ParseFiles() func to template file into a template set. If there is an error
	// we log the detailed error message and use the http.Error() func to send a generic 500 status
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Generic server error
		return
	}

	// We can then use the Execute method on the template set (ts) to write the template content
	// as the response body. The last param to Execute() represents any dynamic data that we want to pass
	// in, which for now will be nil
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err) // Generic server error
	}
}

// Handler for viewing a snippet
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract value of the id param from the query string and attempt to convert it to an integer.
	// If it cannot be converted to an integer, or the value is < 1 then we return a 404
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // Helper for 404s
		return
	}
	// Interpolate the id value with our response and write it to the client
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Handler for creating snippets
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// If the method is not a POST, send a 405 which is a "Method not allowed" status code
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost) // Tell user which methods are available for this endpoint
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
