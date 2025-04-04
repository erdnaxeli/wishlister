package main

import (
	"errors"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// Template store all the availables templates.
//
// It exposes a method Render to render any of the stored templates.
//
// This type is expected to be used with an echo app like this
// ```
// e := echo.New()
// e.Renderer = &Template{...}
// ```
type Template struct {
	listAccessDenied *template.Template
	home             *template.Template
	listEdit         *template.Template
	listNotFound     *template.Template
	listView         *template.Template
	new              *template.Template
}

// ErrUnknownTemplate is returned when we try to render a unknow template.
var ErrUnknownTemplate = errors.New("unknown template")

// Render renders a given template.
//
// If the template is not found, the error ErrUnknownTemplate is returned.
//
// Else, the template is rendered with data and the resulte is written to wr.
func (t *Template) Render(wr io.Writer, name string, data any, _ echo.Context) error {
	switch name {
	case "listAccessDenied":
		return t.listAccessDenied.Execute(wr, data)
	case "home":
		return t.home.Execute(wr, data)
	case "listEdit":
		return t.listEdit.Execute(wr, data)
	case "listNotFound":
		return t.listNotFound.Execute(wr, data)
	case "listView":
		return t.listView.Execute(wr, data)
	case "new":
		return t.new.Execute(wr, data)
	default:
		return ErrUnknownTemplate

	}
}

func loadTemplates(e *echo.Echo) {
	baseTmpl := template.Must(template.ParseFiles("templates/base.html"))

	e.Renderer = &Template{
		listAccessDenied: loadTemplate(baseTmpl, "templates/listAccessDenied.html"),
		home:             loadTemplate(baseTmpl, "templates/index.html"),
		listEdit:         loadTemplate(baseTmpl, "templates/listEdit.html"),
		listNotFound:     loadTemplate(baseTmpl, "templates/listNotFound.html"),
		listView:         loadTemplate(baseTmpl, "templates/listView.html"),
		new:              loadTemplate(baseTmpl, "templates/new.html"),
	}
}

func loadTemplate(baseTmpl *template.Template, filename string) *template.Template {
	return template.Must(template.Must(baseTmpl.Clone()).ParseFiles(filename))
}
