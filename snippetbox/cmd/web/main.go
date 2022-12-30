package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Struct to hold application-wide dependencies for the web app.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	// Define cmd line args
	// Port 4000 will be our default
	addr := flag.String("addr", ":4000", "HTTP Network address")

	// Parses the command line args from the user
	// If we do not call this, it will only use the default argument set by the flag variables
	flag.Parse()

	// Create loggers for appropriate messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// App dependency struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialize our own http.Server struct, so it can use our own pre-defined loggers (above)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Listen on a port and start the server
	// Two parameters are passed in, the TCP network address (port :4000) and the servemux
	infoLog.Printf("Starting server on %s\n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
