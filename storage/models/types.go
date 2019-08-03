package models

import (
	"html/template"
)

type HTMLString string

func (s HTMLString) HTML() template.HTML {
	return template.HTML(s)
}
