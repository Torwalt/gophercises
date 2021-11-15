package internal

import (
	"cyoa/internal/storyteller"
	internalTemplate "cyoa/internal/storyteller/template"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func BuildStory() (*storyteller.Story, error) {
	path, err := filepath.Abs("internal/storyteller/gopher.json")
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path: %v", err)
	}
	storyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not open story file: %v", err)
	}
	story, err := storyteller.New(storyData)
	if err != nil {
		return nil, fmt.Errorf("could not create story from data: %v", err)
	}
	return story, nil
}

func BuildTemplate() (*template.Template, error) {
	path, err := filepath.Abs("internal/storyteller/template/story.html")
	if err != nil {
		return nil, fmt.Errorf("could not get absolute path: %v", err)
	}
	tmpl := internalTemplate.New(path)
	if err != nil {
		return nil, fmt.Errorf("could not create template: %v", err)
	}
	thtml, err := tmpl.AsHTML()
	if err != nil {
		return nil, fmt.Errorf("could not convert template to html: %v", err)
	}

	return thtml, nil
}
