package main

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/zbsss/snippetbox/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}

type templateCache struct {
	cache map[string]*template.Template
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

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/fragments/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		tmplCache.cache[name] = ts
	}

	return &tmplCache, nil
}
