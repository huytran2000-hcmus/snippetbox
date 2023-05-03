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

const (
	flashMessKey = "flash"
	userIDKey    = "authenticatedUserID"
)

type snippetCreateForm struct {
	validator.Validator `form:"-"`
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             string `form:"expires"`
}

type userSignupForm struct {
	validator.Validator `form:"-"`
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
}

type userLoginForm struct {
	validator.Validator `form:"-"`
	Email               string `form:"email"`
	Password            string `form:"password"`
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

	data := app.newDefaultTemplateData(r)
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

	data := app.newDefaultTemplateData(r)
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

	title := form.CheckField("title", form.Title).
		NotBlank("This field can't be blank").
		LE("This field can't be more than 100 characters long", 100).Value()
	content := form.CheckField("content", form.Content).
		NotBlank("This field can't be blank").Value()
	expires := form.CheckField("expires", form.Expires).
		In("This field must be equal 7, 30 or 365", "7", "30", "365").
		ToInt("This field must be a number equal 7, 30 or 365")

	if !form.IsValid() {
		data := app.newDefaultTemplateData(r)
		data.Form = form
		app.render(w, http.StatusBadRequest, "create", data)
		return
	}

	id, err := app.snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), flashMessKey, "Snippet has been successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *Application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	data := app.newDefaultTemplateData(r)
	data.Form = &snippetCreateForm{
		Expires: "365",
	}
	app.render(w, http.StatusOK, "create", data)
}

func (app *Application) userSignupForm(w http.ResponseWriter, r *http.Request) {
	data := app.newDefaultTemplateData(r)
	data.Form = &userSignupForm{}

	app.render(w, http.StatusOK, "signup", data)
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	name := form.CheckField("name", form.Name).
		NotBlank("This field can't be blank").
		LE("This field can't be more than 255 characters long", 255).
		Value()
	email := form.CheckField("email", form.Email).
		NotBlank("This field can't be blank").
		LE("This field can't be more than 255 characters long", 255).
		IsEmail("This field must be a valid email address").
		Value()
	password := form.CheckField("password", form.Password).
		NotBlank("This field can't be blank").
		GE("This field must be at least 8 characters long", 8).
		Value()

	renderFormErrors := func() {
		data := app.newDefaultTemplateData(r)
		data.Form = &form
		app.render(w, http.StatusBadRequest, "signup", data)
	}

	if !form.IsValid() {
		renderFormErrors()
		return
	}

	err = app.user.Insert(name, email, password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "The email address is already in use")
			renderFormErrors()
			return
		}

		if errors.Is(err, models.ErrPasswordTooLong) {
			form.AddFieldError("password", "The password is too long")
			renderFormErrors()
			return
		}

		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), flashMessKey, "You signup was successful. Please log in")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	data := app.newDefaultTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login", data)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	form := userLoginForm{}
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := form.CheckField("email", form.Email).
		NotBlank("This field can't be blank").
		LE("This field can't be more than 255 characters long", 255).
		IsEmail("This field must be a valid email address").
		Value()
	password := form.CheckField("password", form.Password).
		NotBlank("This field can't be blank").
		GE("This field must be at least 8 characters long", 8).
		Value()

	renderFormErrors := func() {
		data := app.newDefaultTemplateData(r)
		data.Form = &form
		app.render(w, http.StatusBadRequest, "login", data)
	}
	if !form.IsValid() {
		renderFormErrors()
		return
	}

	id, err := app.user.Authenticate(email, password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or Password is not correct")
			renderFormErrors()
			return
		}

		app.serverError(w, err)
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), userIDKey, id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}
