package linkparser_test

import (
	"fmt"
	"io"
	"link/pkg/html/linkparser"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		fp       string
		expected []linkparser.Link
	}{
		{
			name: "ex1",
			fp:   "ex1.html",
			expected: []linkparser.Link{
				{Href: "/other-page", Text: "A link to another page"},
			},
		},
		{
			name: "ex2",
			fp:   "ex2.html",
			expected: []linkparser.Link{
				{Href: "https://www.twitter.com/joncalhoun", Text: "Check me out on twitter"},
				{Href: "https://github.com/gophercises", Text: "Gophercises is on Github !"},
			},
		},
		{
			name: "ex3",
			fp:   "ex3.html",
			expected: []linkparser.Link{
				{Href: "#", Text: "Login"},
				{Href: "/lost", Text: "Lost? Need help?"},
				{Href: "https://twitter.com/marcusolsson", Text: "@marcusolsson"},
			},
		},
		{
			name: "ex4",
			fp:   "ex4.html",
			expected: []linkparser.Link{
				{Href: "/dog-cat", Text: "dog cat"},
			},
		},
	}
	for _, test := range tests {

		path, err := filepath.Abs(fmt.Sprintf("testdata/%s", test.fp))
		if err != nil {
			t.Errorf("could not get absolute path: %v", err)
		}

		if err != nil {
			t.Errorf("could not get absolute path: %v", err)
		}
		f, err := os.Open(path)
		if err != nil {
			t.Errorf("could not open file: %v", err)
		}
		defer f.Close()
		var r io.Reader = f
		links, err := linkparser.GetLinks(r)
		assert.Nil(t, err)
		assert.Equal(t, test.expected, links)
	}
}
