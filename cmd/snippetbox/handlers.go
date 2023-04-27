package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newDefaultTemplateData()
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home", data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		app.notFound(w)
		return
	}

	s, err := app.snippet.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	}
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newDefaultTemplateData()
	data.Snippet = s
	app.render(w, http.StatusOK, "view", data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
