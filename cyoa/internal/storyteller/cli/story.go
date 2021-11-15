package cli

import (
	"bufio"
	"cyoa/internal/storyteller"
	"cyoa/internal/storyteller/template"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Handler struct {
	output   io.Writer
	story    storyteller.Story
	template *template.Template
}

func New(output io.Writer, story *storyteller.Story, tmpl *template.Template) *Handler {
	return &Handler{output: output, story: *story, template: tmpl}
}

func (h *Handler) Run() error {
	tmplt, err := h.template.AsText()
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}
	err = tmplt.Execute(os.Stdout, h.story.Intro)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	opts := h.story.Intro.Options

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Choose Path: ")

		a, _ := reader.ReadString('\n')
		a = strings.TrimSpace(a)
		i, err := strconv.Atoi(a)
		if err != nil {
			return fmt.Errorf("could not parse positive integer: %w", err)
		}

		arc := storyteller.Arc{}
		if i > -1 {
			opt := opts[i]
			mArc, ok := h.story.Arcs[opt.Arc]
			arc = mArc
			if !ok {
				return fmt.Errorf("invalid option: %w", err)
			}
		} else {
			arc = h.story.Intro
		}

		opts = arc.Options
		// ideally we would know the width of the screen to fill it with the delim
		delim := strings.Repeat("-", 10)
		fmt.Println(delim)
		err = tmplt.Execute(os.Stdout, arc)
		if err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}
	}

}
