package main

import (
	"net/http"
	"testing"

	"github.com/huytran2000-hcmus/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()
	status, _, body := ts.Get(t, "/ping")

	assert.Equal(t, status, http.StatusOK)

	assert.Equal(t, string(body), "OK")
}

func TestApplication_snippetView(t *testing.T) {
	app := newTestApplication(t)

	srv := newTestServer(t, app.routes())
	defer srv.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := srv.Get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)
			assert.StringContains(t, body, tt.wantBody)
		})
	}
}
