package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	t, ok := app.templates[page]
	if !ok {
		err := fmt.Errorf("the template %q does not exists", page)
		app.serverError(w, err)
		return
	}

	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *Application) newDefaultTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func (app *Application) decodePostForm(r *http.Request, dst interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	var invalidDecoderError *form.InvalidDecoderError
	if errors.As(err, &invalidDecoderError); invalidDecoderError != nil {
		panic(err)
	}

	return err
}
