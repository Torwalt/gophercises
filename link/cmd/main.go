package main

import (
	"flag"
	"fmt"
	"io"
	"link/pkg/html/linkparser"
	"log"
	"os"
	"path/filepath"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("could not run: %v", err)
	}
}

func run() error {
	fp := flag.String("fp", "ex1.html", "file to parse")
	flag.Parse()

	path, err := filepath.Abs(fmt.Sprintf("pkg/html/linkparser/testdata/%s", *fp))
	if err != nil {
		return fmt.Errorf("could not get absolute path: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer f.Close()
	var r io.Reader = f
	links, err := linkparser.GetLinks(r)
	if err != nil {
		return fmt.Errorf("could not get links: %v", err)
	}

	for _, link := range links {
		fmt.Printf("href: %v\n text: %v\n", link.Href, link.Text)
	}

	return nil
}
