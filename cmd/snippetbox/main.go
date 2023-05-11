package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/huytran2000-hcmus/snippetbox/internal/models"
	_ "github.com/lib/pq"
)

type Application struct {
	infoLog        *log.Logger
	errLog         *log.Logger
	snippet        models.Snippets
	users          models.Users
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	debug          bool
}

func main() {
	var addr string
	var dsn string
	var debug bool
	flag.StringVar(&addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&dsn, "dsn", "host=localhost port=5432 user=app_user password=huy2000 dbname=snippetbox sslmode=require search_path=app", "Postgresql datasource name")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR\t", log.LstdFlags|log.Lshortfile)

	db, err := openDB(dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	templates, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.IdleTimeout = 30 * time.Minute
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	app := &Application{
		infoLog:        infoLog,
		errLog:         errLog,
		snippet:        &models.SnippetDB{DB: db},
		users:          &models.UserDB{DB: db},
		templateCache:  templates,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		debug:          debug,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
	}

	srv := http.Server{
		Addr:         addr,
		IdleTimeout:  5 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
	}

	infoLog.Printf("Starting server on %s\n", addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect string is invalid: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping postgresql: %s", err)
	}
	return db, nil
}
