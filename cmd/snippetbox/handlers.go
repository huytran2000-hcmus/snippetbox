package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
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
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.serverError(w, err)
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
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fieldErrs := map[string]string{}
	title = strings.TrimSpace(title)
	if title == "" {
		fieldErrs["title"] = "Field can't be blank"
	}

	if utf8.RuneCountInString(title) > 100 {
		fieldErrs["title"] = "Field can't be more than 100 characters long"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrs["content"] = "Field can't be blank"
	}

	if expires != 7 && expires != 30 && expires != 365 {
		fieldErrs["expires"] = "Field must be equal 7, 30 or 365"
	}

	if len(fieldErrs) > 0 {
		fmt.Fprint(w, fieldErrs)
		return
	}

	id, err := app.snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *Application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	data := app.newDefaultTemplateData()
	app.render(w, http.StatusOK, "create", data)
}
