package template_test

import (
	"cyoa/internal/storyteller/template"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
	sName := "story.html"
	path, err := filepath.Abs(sName)
	if err != nil {
		t.Fatalf("could not get absolute path: %v", err)
	}
	tests := []struct {
		name string
		path string
		err  error
	}{
		{
			name: "success",
			path: path,
			err:  nil,
		},
	}
	for _, test := range tests {
		tmpl := template.New(test.path)
		html, err := tmpl.AsHTML()
		assert.Equal(t, test.err, err)
		assert.Nil(t, err)
		assert.NotNil(t, tmpl)
		assert.Equal(t, html.Name(), sName)
	}
}
