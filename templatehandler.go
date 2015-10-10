package main

import (
	"html/template"
	"sync"
)

const TemplateDir = "templates"

type TemplateHandler struct {
	templ *template.Template
	once  sync.Once
}
