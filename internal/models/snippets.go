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
	Expired time.Time
}

type SnippetRepository struct {
	DB *sql.DB
}

func (m *SnippetRepository) Insert(title string, content string, expired int) (int, error) {
	stmt := "INSERT INTO snippets (title, content, created, expired) values($1, $2, NOW(), NOW() + $3 * INTERVAL '1 DAY') RETURNING id"
	var id int
	err := m.DB.QueryRow(stmt, title, content, expired).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error when inserting snippet: %s", err)
	}

	return id, nil
}

func (m *SnippetRepository) Get(id int) (*Snippet, error) {
	var s Snippet
	stmt := "SELECT id, title, content, created, expired FROM snippets WHERE expired > NOW() AND id = $1"
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expired)

	switch err {
	case sql.ErrNoRows:
		return nil, ErrNoRecord
	case nil:
		return &s, nil
	default:
		return nil, fmt.Errorf("error when selecting a snippet: %s", err)
	}
}

func (m *SnippetRepository) Latest() ([]Snippet, error) {
	stmt := "SELECT id, content, created, expired FROM snippets WHERE expired > NOW() ORDER BY created DESC LIMIT 10"
	row, err := m.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("error when selecting lastest snippets: %s", err)
	}
	defer row.Close()

	var snippets []Snippet
	for row.Next() {
		var s Snippet
		err := row.Scan(&s.ID, &s.Content, &s.Created, &s.Expired)
		if err != nil {
			return nil, fmt.Errorf("error when scaning snippet row: %s", err)
		}
		snippets = append(snippets, s)
	}

	err = row.Err()
	if err != nil {
		return nil, fmt.Errorf("error when iterating snippet row: %s", err)
	}

	return snippets, nil
}
