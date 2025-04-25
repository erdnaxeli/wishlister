// Package statictemplates provides objects to generate go code to render templates.
package statictemplates

// App is the main object to generate code from templates.
//
// It read the templates, validate them and then generate the code.
type App struct {
	generator Generator
	reader    TemplateReader
	validator TemplateValidator
}

// NewWithDefaults returns a new App object with default generator, reader and validator.
func NewWithDefaults(
	templatesDirectory string,
	outputDirectory string,
	packageName string,
) (App, error) {
	generator, err := NewGenerator(outputDirectory, packageName)
	if err != nil {
		return App{}, err
	}

	return New(
		generator,
		NewTemplateReader(templatesDirectory),
		NewTemplateValidator(),
	), nil
}

// New returns a new App object.
func New(generator Generator, reader TemplateReader, validator TemplateValidator) App {
	return App{
		generator: generator,
		reader:    reader,
		validator: validator,
	}
}

// Run reads the templates, validates them and generate the go code.
func (a App) Run() error {
	templates, err := a.reader.ReadTemplates()
	if err != nil {
		return err
	}

	err = a.validator.Validate(templates)
	if err != nil {
		return err
	}

	err = a.generator.Generate(templates)
	if err != nil {
		return err
	}

	return nil
}
