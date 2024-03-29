package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"jeisaRaja.git/snippetbox/pkg/forms"
	"jeisaRaja.git/snippetbox/pkg/models"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	CSRFToken         string
	AuthenticatedUser *models.User
	CurrentYear       int
	Snippet           *models.Snippet
	Snippets          []*models.Snippet
	Form              *forms.Form
	Flash             string
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	fmt.Println(cache)
	return cache, nil
}
