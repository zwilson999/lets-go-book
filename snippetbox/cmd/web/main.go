package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	// Define cmd line args
	addr := flag.String("addr", ":4000", "HTTP Network address")

	// Parses the command line args from the user
	// If we do not call this, it will only use the default argument set by the flag variables
	flag.Parse()

	// Create a new logger variable for info
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create an error logger
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

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
	infoLog.Printf("Starting server on %s\n", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
