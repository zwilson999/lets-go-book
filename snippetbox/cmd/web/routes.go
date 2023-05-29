package main

import "net/http"

// this method initializes our servemux with our routes for the web application
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// pass our servemux as the next parameter to the secureHeaders middleware.
	// because secureHeaders() is just a function, and the function returns a http.Handler, it will
	// also, wrap the mux and secureHeaders with our application infoLog logger
	return app.logRequest(secureHeaders(mux))
}
