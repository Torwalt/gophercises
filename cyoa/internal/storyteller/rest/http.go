package rest

import (
	"cyoa/internal/storyteller"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
)

type Handler struct {
	router chi.Router
	story  storyteller.Story
	tmpl   template.Template
}

func New(router chi.Router, story *storyteller.Story, tmpl *template.Template) *Handler {
	h := &Handler{router: router, story: *story, tmpl: *tmpl}
	h.genRoutes()
	return h
}

func (h *Handler) genRoutes() {
	h.router.Get("/", h.handleRoot)
	h.router.Get("/{arc}", h.handleArc)
}

func (h *Handler) handleRoot(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.Execute(w, h.story.Intro)
	if err != nil {
		renderErr(w, err)
	}
}

func (h *Handler) handleArc(w http.ResponseWriter, r *http.Request) {
	arc, ok := h.story.Arcs[chi.URLParam(r, "arc")]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err := h.tmpl.Execute(w, arc)
	if err != nil {
		renderErr(w, err)
	}
}

func renderErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Errorf("failed to execute template: %v", err).Error()))
}
