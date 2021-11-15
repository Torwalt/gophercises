package main

import (
	"cyoa/cmd/internal"
	"cyoa/internal/storyteller/cli"
	"cyoa/internal/storyteller/template"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error when running the cli: %v", err)
	}
}

func run() error {
	story, err := internal.BuildStory()
	if err != nil {
		return fmt.Errorf("could not build story: %v", err)
	}

	path, err := filepath.Abs("internal/storyteller/template/story.txt")
	if err != nil {
		return fmt.Errorf("could not get absolute path: %v", err)
	}
	tmpl := template.New(path)
	handler := cli.New(os.Stdout, story, tmpl)
	err = handler.Run()
	if err != nil {
		return fmt.Errorf("could not run cli: %v", err)
	}

	return nil
}
