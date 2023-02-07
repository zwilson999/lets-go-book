package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox.lets-go/internal/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox!" as response body
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't then we will send a 404
	if r.URL.Path != "/" {
		app.notFound(w) // Helper for 404s
		return          // Important or else page will keep running
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
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

	data := &templateData{
		Snippets: snippets,
	}

	// We can then use the Execute method on the template set (ts) to write the template content
	// as the response body. The last param to Execute() represents any dynamic data that we want to pass
	// in, which for now will be nil
	err = ts.ExecuteTemplate(w, "base", data)
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

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Initialize a slice containg the paths of the view.tmpl file
	// Plus the base layout and navigation partial
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	// Parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippet: snippet,
	}

	// Execute the template files
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	// Write the snippet data as a plain-text HTTP response body
	fmt.Fprintf(w, "%+v", snippet)
}

// Handler for creating snippets
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// If the method is not a POST, send a 405 which is a "Method not allowed" status code
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // Tell user which methods are available for this endpoint
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Dummy data
	title := "O Snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// Pass the data to SnippetModel.Insert(), receiving the ID of the new record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
