package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type Snippets interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]Snippet, error)
}

type SnippetDB struct {
	DB *sql.DB
}

func (db *SnippetDB) Insert(title string, content string, expires int) (int, error) {
	stmt := "INSERT INTO snippets (title, content, created, expires) values($1, $2, NOW(), NOW() + $3 * INTERVAL '1 DAY') RETURNING id"
	var id int
	err := db.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("models: insert a snippet: %s", err)
	}

	return id, nil
}

func (db *SnippetDB) Get(id int) (*Snippet, error) {
	var s Snippet
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > NOW() AND id = $1"
	err := db.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	switch err {
	case sql.ErrNoRows:
		return nil, ErrNoRecord
	case nil:
		return &s, nil
	default:
		return nil, fmt.Errorf("models: select a snippet: %s", err)
	}
}

func (db *SnippetDB) Latest() ([]Snippet, error) {
	stmt := "SELECT id, title, content, created, expires FROM snippets WHERE expires > NOW() ORDER BY created DESC LIMIT 10"
	row, err := db.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("models: select lastest snippets: %s", err)
	}
	defer row.Close()

	var snippets []Snippet
	for row.Next() {
		var s Snippet
		err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, fmt.Errorf("models: scan snippet row: %s", err)
		}
		snippets = append(snippets, s)
	}

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("models: iterate snippet row: %s", err)
	}

	return snippets, nil
}
