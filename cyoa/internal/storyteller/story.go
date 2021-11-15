package storyteller

import (
	"encoding/json"
)

type Story struct {
	Intro Arc            `json:"intro"`
	Arcs  map[string]Arc `json:"arcs"`
}

type Options struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
type Arc struct {
	Title   string    `json:"title"`
	Story   []string  `json:"story"`
	Options []Options `json:"options"`
}

func New(storyData []byte) (*Story, error) {
	story := Story{}
	if err := json.Unmarshal(storyData, &story); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(storyData, &story.Arcs); err != nil {
		return nil, err
	}
	delete(story.Arcs, "intro")
	return &story, nil
}
