package main

import (
	"cyoa/cmd/internal"
	"cyoa/internal/storyteller/rest"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error when running the server: %v", err)
	}
}

func run() error {
	story, err := internal.BuildStory()
	if err != nil {
		return fmt.Errorf("could not build story: %v", err)
	}

	tmpl, err := internal.BuildTemplate()
	if err != nil {
		return fmt.Errorf("could not build template: %v", err)
	}

	r := chi.NewRouter()
	rest.New(r, story, tmpl)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
	return nil
}
