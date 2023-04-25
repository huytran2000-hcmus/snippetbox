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
		return 0, fmt.Errorf("error when insert snippet: %s", err)
	}

	return id, nil
}

func (m *SnippetRepository) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetRepository) Latest() ([]*Snippet, error) {
	return nil, nil
}
