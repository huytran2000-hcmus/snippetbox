package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
	"github.com/huytran2000-hcmus/snippetbox/ui"
)

type templateData struct {
	Snippet         *models.Snippet
	Snippets        []models.Snippet
	CurrentYear     int
	Form            interface{}
	FlashMessage    string
	IsAuthenticated bool
	CSRFToken       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	files, err := fs.Sub(ui.Files, "html")
	if err != nil {
		return nil, fmt.Errorf("template: don't have html directory: %s", err)
	}
	pages, err := fs.Glob(files, "pages/*")
	if err != nil {
		return nil, fmt.Errorf("template: don't have html pages: %s", err)
	}

	funcMap := template.FuncMap{
		"split_new_line": splitNewLine,
		"timestamp":      timestamp,
		"readable_date":  readableDate,
	}

	for _, page := range pages {
		full_name := filepath.Base(page)
		name := strings.SplitN(full_name, ".", 2)[0]

		t, err := template.New(name).Funcs(funcMap).ParseFS(files, "base.*", "partials/*", page)
		if err != nil {
			return nil, fmt.Errorf("template: error when parsing template files: %s", err)
		}

		cache[name] = t
	}

	return cache, nil
}

func splitNewLine(s string) []string {
	return strings.Split(s, `\n`)
}

func timestamp(t *time.Time) string {
	return t.Format("2006-01-02T15:04:05 -0700 MST")
}

func readableDate(t *time.Time) string {
	return t.Format("Monday, 02 Jan 2006 15:04:05")
}
