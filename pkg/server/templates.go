// Code generated by statictemplates DO NOT EDIT.

package server

import (
	"bytes"
	"html/template"
	"io"
)

type Templates interface {
	RenderBase(wr io.Writer, data any) error
	RenderBaseBytes(data any) ([]byte, error)
	RenderIndex(wr io.Writer, data any) error
	RenderIndexBytes(data any) ([]byte, error)
	RenderListAccessDenied(wr io.Writer, data any) error
	RenderListAccessDeniedBytes(data any) ([]byte, error)
	RenderListEdit(wr io.Writer, data any) error
	RenderListEditBytes(data any) ([]byte, error)
	RenderListNotFound(wr io.Writer, data any) error
	RenderListNotFoundBytes(data any) ([]byte, error)
	RenderListView(wr io.Writer, data any) error
	RenderListViewBytes(data any) ([]byte, error)
	RenderNew(wr io.Writer, data any) error
	RenderNewBytes(data any) ([]byte, error)
	RenderNewGroup(wr io.Writer, data any) error
	RenderNewGroupBytes(data any) ([]byte, error)
	RenderNotFoundError(wr io.Writer, data any) error
	RenderNotFoundErrorBytes(data any) ([]byte, error)
}
type templates struct {
	templateBase             *template.Template
	templateIndex            *template.Template
	templateListAccessDenied *template.Template
	templateListEdit         *template.Template
	templateListNotFound     *template.Template
	templateListView         *template.Template
	templateNew              *template.Template
	templateNewGroup         *template.Template
	templateNotFoundError    *template.Template
}

func NewTemplates() Templates {
	baseTmpl := template.Must(template.New("base.html").Parse("{{ block \"base\" . }}\n<!doctype html>\n<html>\n    <head>\n        <meta charset=\"UTF-8\">\n        <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n        <link href=\"/css/pico.css\" rel=\"stylesheet\">\n        <script defer src=\"https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js\"></script>\n        <script src=\"https://unpkg.com/htmx.org@2.0.4\"></script>\n    </head>\n    <body>\n        <main class=\"container\">\n            <h1>Ma liste de vœux !</h1>\n            {{ block \"content\" . }} Nothing to see here. {{ end }}\n        </main>\n    </body>\n</html>\n{{ end }}\n"))
	notFoundErrorTmpl := template.Must(template.Must(baseTmpl.Clone()).Parse("{{/* base: base.html */}}\n{{ define \"content\" }}\n<p>Page inconnue</p>\n<p><a href=\"/\">Retourner à l'accueil</a></p>\n{{ end }}\n"))
	newGroupTmpl := template.Must(template.Must(baseTmpl.Clone()).Parse("{{/* base: base.html */}}\n{{ define \"content\" }}\n<div>\n  <h2>Créer un groupe</h2>\n  <form method=\"POST\">\n    <fieldset>\n      <label>\n        Nom du groupe\n        <input type=\"text\" name=\"name\" id=\"name\" />\n      </label>\n      <label>\n        Nom d'utilisateur\n        <input type=\"text\" name=\"user\" id=\"user\" />\n      </label>\n      <label>\n        Addresse email. Cela permet de recevoir le lien d'administration du groupe par email et de le retrouver si vous l'avez perdue.\n        <input type=\"text\" name=\"email\" id=\"email\" />\n      </label>\n    </fieldset>\n\n    <input type=\"submit\" value=\"Créer\" />  \n  </form>\n</div>\n{{ end }}\n"))
	newTmpl := template.Must(template.Must(baseTmpl.Clone()).Parse("{{/* base: base.html */}}\n{{ define \"content\" }}\n<div>\n  <h2>Créer une liste de vœux</h2>\n  <form method=\"POST\">\n    <fieldset>\n      <label>\n        Nom d'utilisateur\n        <input type=\"text\" name=\"user\" id=\"user\" />\n      </label>\n      <label>\n        Nom de la liste\n        <input type=\"text\" name=\"name\" id=\"name\" />\n      </label>\n      <label>\n        Addresse email (optionnel). Cela permet de recevoir le lien de la liste de vœux par email et de la retrouver si vous l'avez perdue.\n        <input type=\"text\" name=\"email\" id=\"email\" />\n      </label>\n    </fieldset>\n\n    <input type=\"submit\" value=\"Créer\" />  \n  </form>\n</div>\n{{ end }}\n"))
	listViewTmpl := template.Must(template.Must(baseTmpl.Clone()).Parse("{{/* base: base.html */}}\n{{ define \"content\" }}\n<h2>\n    {{ .Name }}\n    {{ if .AdminID }} (<a href=\"/l/{{ .ID }}/{{ .AdminID }}/edit\">éditer</a>){{ end }}\n</h2>\n\n{{ if .GroupID }}\n  Cette liste fait partie d'un <a href=\"https://www.malistedevoeux.fr/g/{{ .GroupID }}\">groupe</a>.\n{{ end }}\n\n{{ if .AdminID }}\n    <article>\n        <div>Lien à partager : <a href=\"https://malistedevoeux.fr/l/{{ .ID }}\">https://malistedevoeux.fr/l/{{ .ID }}</a></div>\n        <div>Lien d'aministration : <a href=\"https://malistedevoeux.fr/l/{{ .ID }}/{{ .AdminID }}\">https://malistedevoeux.fr/l/{{ .ID }}/{{ .AdminID }}</a></div>\n    </article>\n{{ end }}\n\n<ul>\n    {{ range .Elements }}\n        <li>\n            {{ .Name }}\n            {{ if .Description }} : {{ .Description }}{{ end }}\n            {{ if .URL }} (<a href='{{ .URL }}'>lien</a>){{ end }}\n        </li>\n    {{ end }}\n</ul>\n{{ end }}\n"))
	listNotFoundTmpl := template.Must(template.Must(baseTmpl.Clone()).Parse("{{/* base: base.html */}}\n{{ define \"content\" }}\n<h2>Erreur</h2>\n\n<p>La liste n'a pas pu être trouvée.</p>\n\n<p><a href=\"/\">Retourner à l'accueil</a></p>\n{{ end }}\n"))
	listEditTmpl := template.Must(template.Must(baseTmpl.Clone()).Parse("{{/* base: base.html */}}\n{{ define \"content\" }}\n<h2>Éditer la liste de vœux \"{{ .Name }}\"</h2>\n\n<form method=\"POST\" x-data='{ data: {{ .Data }} }'>\n    <template x-for=\"(obj, index) in data\">\n        <fieldset>\n            <template x-if=\"data[index]['name_error']\" >\n                <label>\n                    Nom\n                    <input type=\"text\" name=\"name\" x-model=\"data[index]['name']\" aria-invalid=\"true\" aria-describedby=\"invalid-helper\" />\n                    <small  id=\"invalid-helper\" x-text=\"data[index]['name_error']\"></small>\n                </label>\n            </template>\n            <template x-if=\"!data[index]['name_error']\" >\n                <label>\n                    Nom\n                    <input type=\"text\" name=\"name\" x-model=\"data[index]['name']\" />\n                </label>\n            </template>\n            </label>\n            <template x-if=\"data[index]['description_error']\">\n                <label>\n                    Description (optionnelle)\n                    <input type=\"text\" name=\"description\" x-model=\"data[index]['description']\" aria-invalid=\"true\" aria-describedby=\"invalid-helper\" />\n                    <small  id=\"invalid-helper\" x-text=\"data[index]['description_error']\"></small>\n                </label>\n            </template>\n            <template x-if=\"!data[index]['description_error']\">\n                <label>\n                    Description (optionnelle)\n                    <input type=\"text\" name=\"description\" x-model=\"data[index]['description']\"/>\n                </label>\n            </template>\n            <template x-if=\"data[index]['url_error']\">\n                <label>\n                    Lien vers l'article (optionnel)\n                    <input type=\"text\" name=\"url\" x-model=\"data[index]['url']\" aria-invalid=\"true\" aria-describedby=\"invalid-helper\" />\n                    <small  id=\"invalid-helper\" x-text=\"data[index]['url_error']\"></small>\n                </label>\n            </template>\n            <template x-if=\"!data[index]['url_error']\">\n                <label>\n                    Lien vers l'article (optionnel)\n                    <input type=\"text\" name=\"url\" x-model=\"data[index]['url']\"/>\n                </label>\n            </template>\n            <button\n                @click=\"data.splice(index, 1)\"\n                type=\"button\"\n            >Supprimer</button>\n            <hr />\n        </fieldset>\n    </template>\n\n    <button\n        @click=\"data.push({ 'id': crypto.randomUUID(), 'name': '', 'description': '', 'url': ''})\"\n        type=\"button\"\n    >Ajouter un nouvel élément</button>\n\n    <input type=\"submit\" value=\"Enregistrer\" />\n</form>\n{{ end }}\n"))
	listAccessDeniedTmpl := template.Must(template.Must(baseTmpl.Clone()).Parse("{{/* base: base.html */}}\n{{ define \"content\" }}\n<div>\n    L'URL est incorrect, vous ne pouvez pas éditer la liste de vœux {{ .Name }}.\n    Vous pouvez cependant <a href=\"/l/{{ .ID }}\">la consulter</a>.\n</div>\n{{ end }}\n"))
	indexTmpl := template.Must(template.Must(baseTmpl.Clone()).Parse("{{/* base: base.html */}}\n{{ define \"content\" }}\n<div>\n    <a href=\"/new\">Créer une nouvelle liste de vœux</a>\n    <p>Cette liste de vœux pourra être consulté par toute personne disposant du lien.</p>\n</div>\n<div>\n    <a href=\"/group/new\">Créer un groupe</a>\n    <p>\n        Vous pourrez inviter des personnes dans le groupe à l'aide de leur adresse email.\n        Chaque membre du groupe pourra créer sa liste de vœux et consulter celle des autres.\n    </p>\n</div>\n{{ end }}\n"))
	return &templates{
		templateBase:             baseTmpl,
		templateIndex:            indexTmpl,
		templateListAccessDenied: listAccessDeniedTmpl,
		templateListEdit:         listEditTmpl,
		templateListNotFound:     listNotFoundTmpl,
		templateListView:         listViewTmpl,
		templateNew:              newTmpl,
		templateNewGroup:         newGroupTmpl,
		templateNotFoundError:    notFoundErrorTmpl,
	}
}
func (t *templates) RenderBase(wr io.Writer, data any) error {
	return t.templateBase.Execute(wr, data)
}
func (t *templates) RenderBaseBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderBase(wr, data)
	return wr.Bytes(), err
}
func (t *templates) RenderIndex(wr io.Writer, data any) error {
	return t.templateIndex.Execute(wr, data)
}
func (t *templates) RenderIndexBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderIndex(wr, data)
	return wr.Bytes(), err
}
func (t *templates) RenderListAccessDenied(wr io.Writer, data any) error {
	return t.templateListAccessDenied.Execute(wr, data)
}
func (t *templates) RenderListAccessDeniedBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderListAccessDenied(wr, data)
	return wr.Bytes(), err
}
func (t *templates) RenderListEdit(wr io.Writer, data any) error {
	return t.templateListEdit.Execute(wr, data)
}
func (t *templates) RenderListEditBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderListEdit(wr, data)
	return wr.Bytes(), err
}
func (t *templates) RenderListNotFound(wr io.Writer, data any) error {
	return t.templateListNotFound.Execute(wr, data)
}
func (t *templates) RenderListNotFoundBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderListNotFound(wr, data)
	return wr.Bytes(), err
}
func (t *templates) RenderListView(wr io.Writer, data any) error {
	return t.templateListView.Execute(wr, data)
}
func (t *templates) RenderListViewBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderListView(wr, data)
	return wr.Bytes(), err
}
func (t *templates) RenderNew(wr io.Writer, data any) error {
	return t.templateNew.Execute(wr, data)
}
func (t *templates) RenderNewBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderNew(wr, data)
	return wr.Bytes(), err
}
func (t *templates) RenderNewGroup(wr io.Writer, data any) error {
	return t.templateNewGroup.Execute(wr, data)
}
func (t *templates) RenderNewGroupBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderNewGroup(wr, data)
	return wr.Bytes(), err
}
func (t *templates) RenderNotFoundError(wr io.Writer, data any) error {
	return t.templateNotFoundError.Execute(wr, data)
}
func (t *templates) RenderNotFoundErrorBytes(data any) ([]byte, error) {
	wr := &bytes.Buffer{}
	err := t.RenderNotFoundError(wr, data)
	return wr.Bytes(), err
}
