package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/zbsss/snippetbox/internal/models"
	"github.com/zbsss/snippetbox/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Toast           string
	IsAuthenticated bool
	CSRFToken       string
}

type templateCache struct {
	cache map[string]*template.Template
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2023 at 15:04")
}

var funcs = template.FuncMap{
	"humanDate": humanDate,
}

func (tc *templateCache) Get(templateName string) (*template.Template, error) {
	ts, ok := tc.cache[templateName]
	if !ok {
		return nil, fmt.Errorf("template %s not found in cache", templateName)
	}

	return ts, nil
}

func newTemplateCache() (*templateCache, error) {
	tmplCache := templateCache{
		cache: make(map[string]*template.Template),
	}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/fragments/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(funcs).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		tmplCache.cache[name] = ts
	}

	return &tmplCache, nil
}
