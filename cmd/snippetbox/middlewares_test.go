package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/huytran2000-hcmus/snippetbox/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	wantHeader := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com; frame-ancestors 'none'"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), wantHeader)

	wantHeader = "strict-origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), wantHeader)

	wantHeader = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), wantHeader)

	wantHeader = "DENY"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), wantHeader)

	wantHeader = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), wantHeader)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	body, err := io.ReadAll(rs.Body)
	defer rs.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(body), "OK")
}
