package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox.lets-go/internal/models"
)

// define a home handler function which writes a byte slice containing
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// check if the current request URL path exactly matches "/". If it doesn't then we will send a 404
	if r.URL.Path != "/" {
		app.notFound(w) // helper for 404s
		return          // important or else page will keep running
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// slice to contain our template files
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	// use the template.ParseFiles() func to template file into a template set. If there is an error
	// we log the detailed error message and use the http.Error() func to send a generic 500 status
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err) // Generic server error
		return
	}

	data := &templateData{
		Snippets: snippets,
	}

	// we can then use the Execute method on the template set (ts) to write the template content
	// as the response body. The last param to Execute() represents any dynamic data that we want to pass
	// in, which for now will be nil
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err) // generic server error
	}
}

// handler for viewing a snippet
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	// extract value of the id param from the query string and attempt to convert it to an integer.
	// if it cannot be converted to an integer, or the value is < 1 then we return a 404
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w) // helper for 404s
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

	// initialize a slice containg the paths of the view.tmpl file
	// plus the base layout and navigation partial
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/view.tmpl.html",
	}

	// parse the template files
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippet: snippet,
	}

	// execute the template files
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	// write the snippet data as a plain-text HTTP response body
	fmt.Fprintf(w, "%+v", snippet)
}

// handler for creating snippets
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	// if the method is not a POST, send a 405 which is a "Method not allowed" status code
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost) // tell user which methods are available for this endpoint
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// dummy data
	title := "O Snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// pass the data to SnippetModel.Insert(), receiving the ID of the new record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
