package storyteller_test

import (
	"cyoa/internal/storyteller"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()
	path, err := filepath.Abs("gopher.json")
	if err != nil {
		t.Fatalf("could not get absolute path: %v", err)
	}
	storyData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read file: %v", err)
	}
	tests := []struct {
		name string
		data []byte
		err  error
	}{
		{
			name: "success",
			data: storyData,
			err:  nil,
		},
	}
	for _, test := range tests {
		story, err := storyteller.New(test.data)
		assert.Equal(t, test.err, err)
		assert.NotNil(t, story)
		_, hasDuplicate := story.Arcs["intro"]
		assert.False(t, hasDuplicate)
	}

}
