package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
)

type templateData struct {
	Snippet     *models.Snippet
	Snippets    []models.Snippet
	CurrentYear int
	Form        interface{}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		"split_new_line": splitNewLine,
		"timestamp":      timestamp,
		"readable_date":  readableDate,
	}
	for _, page := range pages {
		full_name := filepath.Base(page)
		name := strings.SplitN(full_name, ".", 2)[0]

		t, err := template.New(name).Funcs(funcMap).ParseFiles("./ui/html/base.tmpl.html", page)
		if err != nil {
			return nil, fmt.Errorf("error when parsing template files: %s", err)
		}
		t, err = t.ParseGlob("./ui/html/partials/*tmpl.html")
		if err != nil {
			return nil, fmt.Errorf("error when parsing template files: %s", err)
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
