package main

import (
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func newTestApplication() *Application {
	return &Application{
		infoLog: log.New(io.Discard, "", 0),
		errLog:  log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, mux http.Handler) *testServer {
	ts := httptest.NewTLSServer(mux)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	return &testServer{ts}
}

func (ts *testServer) Get(t *testing.T, path string) (int, http.Header, string) {
	t.Helper()

	rs, err := ts.Client().Get(ts.URL + path)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	return rs.StatusCode, rs.Header, string(body)
}
