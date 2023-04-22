package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Controller struct {
	infoLog *log.Logger
	errLog  *log.Logger
}

func (c *Controller) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		c.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/pages/home.gohtml",
		"./ui/html/partials/nav.gohtml",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		c.serverError(w, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		c.serverError(w, err)
	}
}

func (c *Controller) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		c.notFound(w)
		return
	}
	message := fmt.Sprintf("Display snippet with ID %d", id)
	w.Write([]byte(message))
}

func (c *Controller) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		c.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a specific snippet..."))
}
