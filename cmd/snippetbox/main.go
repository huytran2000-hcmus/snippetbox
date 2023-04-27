package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
	_ "github.com/lib/pq"
)

type Application struct {
	infoLog   *log.Logger
	errLog    *log.Logger
	snippet   *models.SnippetRepository
	templates map[string]*template.Template
}

func main() {
	var addr string
	var dsn string
	flag.StringVar(&addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&dsn, "dsn", "host=localhost port=5432 user=app_user password=huy2000 dbname=snippetbox sslmode=require search_path=app", "Postgresql datasource name")
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
	app := &Application{
		infoLog:   infoLog,
		errLog:    errLog,
		snippet:   &models.SnippetRepository{DB: db},
		templates: templates,
	}

	srv := http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     errLog,
		Handler:      app.routes(),
	}
	infoLog.Printf("Starting server on %s\n", addr)
	err = srv.ListenAndServe()
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
