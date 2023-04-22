package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (c *Controller) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	c.errLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (c *Controller) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (c *Controller) notFound(w http.ResponseWriter) {
	c.clientError(w, http.StatusNotFound)
}
