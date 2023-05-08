package main

import (
	"net/http"

	"github.com/huytran2000-hcmus/snippetbox/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	statefulMW := alice.New(app.sessionManager.LoadAndSave, CSRFPrevent, app.authenticate)
	router.Handler(http.MethodGet, "/", statefulMW.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", statefulMW.ThenFunc(app.snippetView))

	router.Handler(http.MethodGet, "/user/signup", statefulMW.ThenFunc(app.userSignupForm))
	router.Handler(http.MethodPost, "/user/signup", statefulMW.ThenFunc(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", statefulMW.ThenFunc(app.userLoginForm))
	router.Handler(http.MethodPost, "/user/login", statefulMW.ThenFunc(app.userLogin))

	protectedMW := statefulMW.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protectedMW.ThenFunc(app.snippetCreateForm))
	router.Handler(http.MethodPost, "/snippet/create", protectedMW.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/user/logout", protectedMW.ThenFunc(app.userLogout))

	standardMW := alice.New(app.recoverFromPanic, app.logRequest, secureHeaders)
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	return standardMW.Then(router)
}
