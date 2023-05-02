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
	validator.Validator `form:"-"`
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             string `form:"expires"`
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
	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := form.Check("title", form.Title).
		NotBlank("This field can't be blank").
		LE("This field can't be more than 100 characters long", 100).FieldValue
	content := form.Check("content", form.Content).
		NotBlank("This field can't be blank").FieldValue
	expires, err := form.Check("expires", form.Expires).
		In("This field must be equal 7, 30 or 365", "7", "30", "365").
		ToInt()
	if err != nil {
		app.errLog.Print(err)
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if !form.IsValid() {
		data := app.newDefaultTemplateData()
		data.Form = form
		app.render(w, http.StatusBadRequest, "create", data)
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
	data.Form = &snippetCreateForm{
		Expires: "365",
	}
	app.render(w, http.StatusOK, "create", data)
}
