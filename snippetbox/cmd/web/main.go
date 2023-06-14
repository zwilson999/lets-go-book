package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"snippetbox.lets-go/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// struct to hold application-wide dependencies
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// use the scs.New() func to init a new session manager
	// then we configure it to use our MySql database as the session store, and set a lifetime of 12 hours
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// app dependency struct
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// initialize a tls.Config struct to hold non-default TLS settings we want our server to use.
	// in this case the only thing we are changing is the curve preferences value, so that the only
	// elliptic curves with assembly implementations are used
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// initialize our own http.Server struct, so it can use our own pre-defined loggers (above)
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// listen on a port and start the server
	// two parameters are passed in, the TCP network address (port :4000) and the servemux
	infoLog.Printf("starting server on %s\n", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
