package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR\t", log.LstdFlags|log.Lshortfile)

	ctrl := &Controller{
		infoLog: infoLog,
		errLog:  errLog,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", ctrl.home)
	mux.HandleFunc("/snippet/view", ctrl.snippetView)
	mux.HandleFunc("/snippet/create", ctrl.snippetCreate)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		ErrorLog:     errLog,
		Handler:      mux,
	}
	infoLog.Printf("Starting server on %s\n", addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}
