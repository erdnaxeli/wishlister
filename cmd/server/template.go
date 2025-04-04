package main

import (
	"errors"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	listAccessDenied *template.Template
	home             *template.Template
	listEdit         *template.Template
	listNotFound     *template.Template
	listView         *template.Template
	new              *template.Template
}

var ErrUnknownTemplate = errors.New("unknown template")

func (t *Template) Render(wr io.Writer, name string, data any, c echo.Context) error {
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
