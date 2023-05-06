package mock

import (
	"time"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
)

var mockSnippet = &models.Snippet{
	ID:      0,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type StubSnippets struct{}

func (s *StubSnippets) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (s *StubSnippets) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (s *StubSnippets) Latest() ([]models.Snippet, error) {
	return []models.Snippet{*mockSnippet}, nil
}
