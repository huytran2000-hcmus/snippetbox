package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("postgres", "host=localhost port=5432 user=test_readwrite password=testsnippetbox dbname=test_snippetbox sslmode=require search_path=app")
	if err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile("testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		script, err := os.ReadFile("testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})

	return db
}
