package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	urlshortener "urlshort/pkg"

	"github.com/boltdb/bolt"
)

func main() {
	// paths from files
	yamlPath := flag.String("yaml", "default-paths.yaml", "yaml file")
	jsonPath := flag.String("json", "default-paths.json", "json file")
	flag.Parse()
	yaml, err := os.ReadFile(*yamlPath)
	if err != nil {
		log.Fatal(err)
	}
	json, err := os.ReadFile(*jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	bName := "paths"
	db, err := initDB(bName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshortener.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshortener.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		log.Fatal(err)
	}

	jsonHandler, err := urlshortener.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		log.Fatal(err)
	}

	boltDBHandler, err := urlshortener.BoltDBHandler(db, jsonHandler, bName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltDBHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func initDB(bName string) (*bolt.DB, error) {
	dbName := "path.db"
	os.Remove(dbName)
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = urlshortener.InitBoltDB(db, "paths")
	if err != nil {
		return nil, err
	}

	data := []urlshortener.PathURL{
		{Path: "/so", URL: "https://stackoverflow.com/"},
		{Path: "/hn", URL: "https://news.ycombinator.com/"},
	}
	err = urlshortener.PopulateBoltDB(db, data, bName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
