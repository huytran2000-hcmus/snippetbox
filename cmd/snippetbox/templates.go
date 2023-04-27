package main

import "github.com/huytran2000-hcmus/snippetbox/internal/models"

type templateData struct {
	Snippet  *models.Snippet
	Snippets []models.Snippet
}
