package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	statefulMW := alice.New(app.sessionManager.LoadAndSave)
	router.Handler(http.MethodGet, "/", statefulMW.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", statefulMW.ThenFunc(app.snippetView))
	router.Handler(http.MethodPost, "/snippet/create", statefulMW.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodGet, "/snippet/create", statefulMW.ThenFunc(app.snippetCreateForm))

	router.Handler(http.MethodGet, "/user/signup", statefulMW.ThenFunc(app.userSignupForm))
	router.Handler(http.MethodPost, "/user/signup", statefulMW.ThenFunc(app.userSignup))

	standardMW := alice.New(app.recoverFromPanic, app.logRequest, secureHeaders)
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	return standardMW.Then(router)
}
