package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.lets-go/internal/models"
)

// acts as structure to hold dynamic data that we want to pass to HTML templates
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string // for holding string data to flash to user once upon certain request
}

// func to format date in a human-readable form
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
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

	// use filepath.Glob() to get a slice of all filepaths that match the pattern .tmpl.html
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {

		// extract the file name from the full filepath and assign it to the name variable
		name := filepath.Base(page)

		// parse the base template into a template set
		// the FuncMap must be registered with the template set before calling ParseFiles()
		// template.New(name) will create an empty template set with the given name
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// use ParseGlob() on the partials template set to add any templates that exist therein
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// parse the template files into a set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template set to map cache
		cache[name] = ts
	}
	return cache, nil
}
