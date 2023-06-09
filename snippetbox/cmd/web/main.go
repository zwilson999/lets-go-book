package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	"snippetbox.lets-go/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// struct to hold application-wide dependencies
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {

	// define cmd line args
	addr := flag.String("addr", ":4000", "HTTP Network address")
	dsn := flag.String("dsn", "", "MySql Data Source Name. should be in the form web:pass@/snippetbox?parseTime=true")

	// parses the command line args from the user
	// if we do not call this, it will only use the default argument set by the flag variables
	flag.Parse()

	// create loggers for appropriate messages
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// connect to DB
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// init new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// init form decoder
	formDecoder := form.NewDecoder()

	// app dependency struct
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// initialize our own http.Server struct, so it can use our own pre-defined loggers (above)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// listen on a port and start the server
	// two parameters are passed in, the TCP network address (port :4000) and the servemux
	infoLog.Printf("starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// function to open mysql db and return pointer to db handle
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
