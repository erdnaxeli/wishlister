package statictemplates

import (
	"fmt"
	htmlTemplate "html/template"
	textTemplate "text/template"
)

// TemplateValidator is the main interface to validate some templates.
//
// Validating a template means that it can correctly be parsed by text/template or
// html/template.
type TemplateValidator interface {
	// Validate returns nil if all the templates can be parsed without error.
	//
	// If there is any error, the error is returned.
	//
	// If a template use another one as its base, both templates are parsed in the same
	// template object. If the base template also has a base, it is parsed with its own
	// base, and so on. Loops are not prevented, it will end with a stack overflow error.
	//
	// If a template use another one as its base, but the other one is missing, an error
	// is returned.
	//
	// If the name of a template ends with ".html", it is parsed using html/template,
	// else using text/template. The base template will be parsed using the same module,
	// regardless its actual name.
	Validate(templates map[string]Template) error
}

type templateValidator struct{}

// NewTemplateValidator return a new TemplateValidator object.
func NewTemplateValidator() TemplateValidator {
	return templateValidator{}
}

func (tv templateValidator) Validate(templates map[string]Template) error {
	ctv := cacheTemplateValidator{
		templates: templates,
		tmplHTML:  make(map[string]*htmlTemplate.Template),
		tmplText:  make(map[string]*textTemplate.Template),
	}

	return ctv.validate()
}

type cacheTemplateValidator struct {
	templates map[string]Template

	// caches for already parsed templates
	tmplHTML map[string]*htmlTemplate.Template
	tmplText map[string]*textTemplate.Template
}

func (ctv *cacheTemplateValidator) validate() error {
	for name, template := range ctv.templates {
		err := ctv.validateOne(name, template)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ctv *cacheTemplateValidator) validateOne(name string, template Template) error {
	var err error

	if name[len(name)-5:] == ".html" {
		_, err = ctv.validateOneHTML(name, template)
	} else {
		_, err = ctv.validateOneText(name, template)
	}

	return err
}

func (ctv *cacheTemplateValidator) validateOneHTML(
	name string,
	template Template,
) (*htmlTemplate.Template, error) {
	var err error
	var tmpl *htmlTemplate.Template

	if template.Base == "" {
		tmpl, err = htmlTemplate.New(name).Parse(template.Content)
	} else {
		baseTmpl, ok := ctv.tmplHTML[template.Base]
		if !ok {
			baseTemplate, ok := ctv.templates[template.Base]
			if !ok {
				return nil, BaseNotFoundError{
					base:     template.Base,
					template: name,
				}
			}

			baseTmpl, err = ctv.validateOneHTML(template.Base, baseTemplate)
			if err != nil {
				return nil, err
			}
		}

		tmpl, err = baseTmpl.Parse(template.Content)
	}

	if err != nil {
		return nil, fmt.Errorf("error while parsing %s: %w", name, err)
	}

	return tmpl, nil
}

func (ctv *cacheTemplateValidator) validateOneText(
	name string,
	template Template,
) (*textTemplate.Template, error) {
	var err error
	var tmpl *textTemplate.Template

	if template.Base == "" {
		tmpl, err = textTemplate.New(name).Parse(template.Content)
	} else {
		baseTmpl, ok := ctv.tmplText[template.Base]
		if !ok {
			baseTemplate, ok := ctv.templates[template.Base]
			if !ok {
				return nil, BaseNotFoundError{
					base:     template.Base,
					template: name,
				}
			}

			baseTmpl, err = ctv.validateOneText(template.Base, baseTemplate)
			if err != nil {
				return nil, err
			}
		}

		tmpl, err = baseTmpl.Parse(template.Content)
	}

	if err != nil {
		return nil, fmt.Errorf("error while parsing %s: %w", name, err)
	}

	return tmpl, nil
}
