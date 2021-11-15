package template

import (
	"fmt"
	"html/template"
	"strings"
	txtTemplate "text/template"
)

type Template struct {
	Path string
}

func New(path string) *Template {
	return &Template{Path: path}
}

func (t *Template) AsHTML() (*template.Template, error) {
	tmpl := template.New("story.html")
	funcMap := template.FuncMap{
		"Titleize": strings.Title,
	}
	tmpl = tmpl.Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(t.Path)
	if err != nil {
		return nil, fmt.Errorf("template.ParseFiles: %w", err)
	}
	return tmpl, nil
}

func (t *Template) AsText() (*txtTemplate.Template, error) {
	tmpl := txtTemplate.New("story.txt")
	funcMap := txtTemplate.FuncMap{
		"Titleize": strings.Title,
	}
	tmpl = tmpl.Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(t.Path)
	if err != nil {
		return nil, fmt.Errorf("template.ParseFiles: %w", err)
	}
	return tmpl, nil
}

// func (t *Template) createbaseTemplate() (*txtTemplate.Template, error) {

// }
