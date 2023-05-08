package main

import (
	"net/http"
	"net/url"
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

func TestUserSignUp(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.Get(t, "/user/signup")
	token := extractCSRFToken(t, body)

	const (
		name     = "Huy"
		email    = "abc@gmail.com"
		password = "password"
		formTag  = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct {
		name         string
		userName     string
		userPassword string
		userEmail    string
		csrfToken    string
		wantStatus   int
		formTag      string
	}{
		{
			name:         "Valid Submission",
			userName:     name,
			userPassword: password,
			userEmail:    email,
			csrfToken:    token,
			wantStatus:   http.StatusSeeOther,
		},
		{
			name:         "Missing csqf token",
			userName:     name,
			userPassword: password,
			userEmail:    email,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "Empty Name",
			userName:     "",
			userPassword: password,
			userEmail:    email,
			csrfToken:    token,
			wantStatus:   http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
		{
			name:         "Invalid email",
			userName:     name,
			userPassword: password,
			userEmail:    "bob@email.",
			csrfToken:    token,
			wantStatus:   http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
		{
			name:         "Email already in use",
			userName:     name,
			userPassword: password,
			userEmail:    "dupe@example.com",
			csrfToken:    token,
			wantStatus:   http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
		{
			name:         "Short password",
			userName:     name,
			userPassword: "123",
			userEmail:    email,
			csrfToken:    token,
			wantStatus:   http.StatusUnprocessableEntity,
			formTag:      formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var form url.Values = map[string][]string{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)
			status, _, body := ts.PostForm(t, "/user/signup", form)

			assert.Equal(t, status, tt.wantStatus)
			assert.StringContains(t, body, tt.formTag)
		})
	}
}
