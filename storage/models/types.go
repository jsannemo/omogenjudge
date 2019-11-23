package models

import (
	"html/template"
)

// An HTMLString represents a string that is assumed to be safe HTML.
type HTMLString string

// HTML converts a HTMLString to a template variable for HTML.
func (s HTMLString) HTML() template.HTML {
	return template.HTML(s)
}
