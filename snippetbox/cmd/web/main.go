package main

import (
	"log"
	"net/http"
)

func main() {
	// Createa servemux and register the home function as the handler for the root pattern
	mux := http.NewServeMux()

	// Create a file server which serves files out of our static directory
	// Will be relative to root of the directory
	// Note we will also strip the static prefix to search for files within the directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Page routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Listen on a port and start the server
	// Two parameters are passed in, the TCP network address (port :4000)
	// and the servemux
	log.Print("Starting server on port:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
