package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.lets-go/ui"

	"snippetbox.lets-go/internal/models"
)

// acts as structure to hold dynamic data that we want to pass to HTML templates
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string // for holding string data to flash to user once upon certain request
	IsAuthenticated bool
	CSRFToken       string
}

// func to format date in a human-readable form
func humanDate(t time.Time) string {

	// return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// initialize a FuncMap and store it as a global variable.
// this is basically a string-keyed map which acts as a lookup between the names
// of our custom template functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	// initialize a new map to act as our cache
	cache := map[string]*template.Template{}

	// use the fs.Glob() to get a slice of all the filepaths in the ui.Files embedded filesystem
	// which match the pattern 'html/pages/*.tmpl.html'. This essentially gives us a slice of all the page templates for the app,
	// just like before.

	// use filepath.Glob() to get a slice of all filepaths that match the pattern .tmpl.html
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		// extract the file name from the full filepath and assign it to the name variable
		name := filepath.Base(page)

		// create a slice containing the filepath patterns for the templates we want to parse
		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		// parse the base template into a template set
		// template.New(name) will create an empty template set with the given name
		// use ParseFS() instead of ParseFiles() to parse the template files from the ui.files embedded FS
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// add template set to map cache
		cache[name] = ts
	}
	return cache, nil
}
