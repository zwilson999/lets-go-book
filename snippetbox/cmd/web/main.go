package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"snippetbox.lets-go/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Struct to hold application-wide dependencies for the web app.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {

	// Define cmd line args
	// Port 4000 will be our default
	addr := flag.String("addr", ":4000", "HTTP Network address")

	// DSN for connecting to MySQL
	dsn := flag.String("dsn", "web:Bramble187*@/snippetbox?parseTime=true", "MySql Data Source Name")

	// Parses the command line args from the user
	// If we do not call this, it will only use the default argument set by the flag variables
	flag.Parse()

	// Create loggers for appropriate messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Connect to DB
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Close connection pool before main() exits
	defer db.Close()

	// App dependency struct
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
