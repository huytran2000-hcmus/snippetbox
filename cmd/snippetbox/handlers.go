package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
	"github.com/huytran2000-hcmus/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	validator.Validator
	Title   string
	Content string
	Expires string
}

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

	titleRaw := r.PostForm.Get("title")
	contentRaw := r.PostForm.Get("content")
	expiresRaw := r.PostForm.Get("expires")
	form := snippetCreateForm{
		Title:   titleRaw,
		Content: contentRaw,
		Expires: expiresRaw,
	}

	form.Check("title", titleRaw).
		NotBlank("This field can't be blank").
		MaxCharacters("This field can't be more than 100 characters long", 100)
	form.Check("content", contentRaw).
		NotBlank("This field can't be blank")
	form.Check("expires", expiresRaw).
		InPermittedArr("This field must be equal 7, 30 or 365", "7", "30", "365")
	if !form.IsValid() {
		data := app.newDefaultTemplateData()
		data.Form = form
		app.render(w, http.StatusBadRequest, "create", data)
		return
	}

	title := titleRaw
	content := contentRaw
	expires, err := strconv.Atoi(expiresRaw)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
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
	data.Form = &snippetCreateForm{
		Expires: "365",
	}
	app.render(w, http.StatusOK, "create", data)
}
