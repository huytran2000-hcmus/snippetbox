package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/huytran2000-hcmus/snippetbox/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		"split_new_line": splitNewLine,
	}
	for _, page := range pages {
		full_name := filepath.Base(page)
		name := strings.TrimSuffix(full_name, ".tmpl.html")

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
