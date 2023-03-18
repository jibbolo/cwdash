package widget

import (
	"html/template"

	"github.com/russross/blackfriday/v2"
)

type Widget struct {
	Properties map[string]any `json:"properties,omitempty"`
}

func (w Widget) HasMarkdown() bool {
	_, ok := w.Properties["markdown"]
	return ok
}

func (w Widget) HasQuery() bool {
	_, ok := w.Properties["query"]
	return ok
}

func (w *Widget) RenderMarkdown() template.HTML {
	if markdown, ok := w.Properties["markdown"]; ok {
		return template.HTML(
			blackfriday.Run([]byte(markdown.(string))),
		)
	}
	return "Error: markdown unavailable"
}
