package main

import (
	"net/http"
	"testing"

	"github.com/huytran2000-hcmus/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication()

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	status, _, body := ts.Get(t, "/ping")

	assert.Equal(t, status, http.StatusOK)

	assert.Equal(t, string(body), "OK")
}
