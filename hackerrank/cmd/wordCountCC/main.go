package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	s := flag.String("s", "TestCase", "camelcase string to count words")
	flag.Parse()

	n, err := countWordsCamelCase(*s)
	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf("%s has %d words\n", *s, n)
}

func countWordsCamelCase(s string) (int, error) {
	if len(s) < 1 {
		return 0, fmt.Errorf("%s is too short", s)
	}

	wordCount := 1
	for pos, char := range s {
		if pos == 0 {
			continue
		}

		// if char is upper case, increment word count
		if string(char) == strings.ToUpper(string(char)) {
			wordCount++
		}
	}

	return wordCount, nil
}
