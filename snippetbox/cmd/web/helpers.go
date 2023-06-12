package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

// this serverError helper writes an error message and stack trace to errorLog
// then sends a generic 500 Internal Server Error response to user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack()) // debug.Stack() gets stack trace for current goroutines
	app.errorLog.Output(2, trace)                              // change depth of stack trace

	// write internal server error to response writer
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// this helper sends specific status codes to the user as well as the status text
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// this helper constructs a clientError wrapper for statusNotFound which will send a 404
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	// retrieve the appropriate template set from our app cache. if no such entry exists, then create a new error
	// and call the serverError() helper
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// create new trial buffer
	buf := new(bytes.Buffer)

	// write template to the trial buffer to test that our template write works, instead of straight to the writer. if there is an error
	// call our serverError() helper
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// write the provided status
	w.WriteHeader(status)

	// write contents of buffer to the writer.
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// add flash message to template data if one exists
		// this will be triggered to user when they create a snippet. otherwise it will be an empty string and will
		// not be rendered in the template display
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
	}
}

// create new decodePostForm() helper method. the second param "dst" is the target destination
// that we want to decode the form data into
func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// call form Decode()
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// if we try to use an invalid target destination, the Decode() method
		// will return an error of the type *form.InvalidDecoderError
		// we use errors.As() to check for this specific error and panic rather than returning the rror
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// return all other errors
		return err
	}
	return nil
}
