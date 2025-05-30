package statictemplates

import (
	"fmt"
	"io"
	"iter"
	"maps"
	"os"
	"path"
	"slices"

	"github.com/dave/jennifer/jen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Generator is the main interface to generate go code.
type Generator interface {
	// Generate generates the go code corresponding to the given templates.
	//
	// It generates a single file "templates.go".
	Generate(templates map[string]Template) error
}

type generator struct {
	caser       cases.Caser
	packageName string
	writer      DirectoryWriter
}

// NewGenerator return a object implementing the Generator interface.
//
// All code generated using this object will use html/template.
//
// The code will be written in the given directory.
func NewGenerator(directory string, packageName string) (Generator, error) {
	writer, err := NewOsDirectoryWriter(directory)
	if err != nil {
		return nil, err
	}

	return NewGeneratorWriter(writer, packageName), nil
}

// NewGeneratorWriter return a object implementing the Generator interface.
//
// All code generated using this object will use html/template.
//
// The code will be written using the given DirectoryWriter object. Unless you have
// very specific needs, you should use NewGenerator instead.
func NewGeneratorWriter(writer DirectoryWriter, packageName string) Generator {
	return generator{
		caser:       cases.Title(language.English, cases.NoLower),
		packageName: packageName,
		writer:      writer,
	}
}

// Generate generates the go code and write it to a file.
func (g generator) Generate(
	templates map[string]Template,
) error {
	file := jen.NewFile(g.packageName)
	file.HeaderComment("Code generated by statictemplates DO NOT EDIT.")

	templatesIter := MapOrdered(templates)
	g.generateInterface(file, templatesIter)
	g.generateStruct(file, templatesIter)
	err := g.generateStructConstructor(file, templates)
	if err != nil {
		return err
	}

	g.generateStructMethods(file, templatesIter)

	f, err := g.writer.Create("templates.go")
	if err != nil {
		return fmt.Errorf("error while creating file templates.go: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	err = file.Render(f)
	if err != nil {
		return fmt.Errorf("error while writing file templates.go: %w", err)
	}

	return nil
}

func (g generator) generateInterface(file *jen.File, templates iter.Seq2[string, Template]) {
	methods := make([]jen.Code, 0)

	for name := range templates {
		name = g.getName(name)
		methodName := fmt.Sprintf("Render%s", g.caser.String(name))
		methodBytesName := fmt.Sprintf("Render%sBytes", g.caser.String(name))
		method := jen.Id(methodName).Params(
			jen.Id("wr").Qual("io", "Writer"),
			jen.Id("data").Any(),
		).Error()
		methodBytes := jen.Id(methodBytesName).
			Params(jen.Id("data").Any()).
			Params(jen.Index().Byte(), jen.Error())
		methods = append(methods, method, methodBytes)
	}

	file.Type().Id("Templates").Interface(methods...)
}

func (g generator) generateStruct(file *jen.File, templates iter.Seq2[string, Template]) {
	fields := make([]jen.Code, 0)

	for name := range templates {
		name = g.getName(name)
		fieldName := fmt.Sprintf("template%s", g.caser.String(name))
		field := jen.Id(fieldName).Add(jen.Op("*")).Qual("html/template", "Template")

		fields = append(fields, field)
	}

	file.Type().Id("templates").Struct(fields...)
}

func (g generator) generateStructMethods(file *jen.File, templates iter.Seq2[string, Template]) {
	for name := range templates {
		name := g.getName(name)
		fieldName := fmt.Sprintf("template%s", g.caser.String(name))
		methodName := fmt.Sprintf("Render%s", g.caser.String(name))
		params := []jen.Code{jen.Id("wr").Qual("io", "Writer"), jen.Id("data").Any()}
		methodBytesName := fmt.Sprintf("Render%sBytes", g.caser.String(name))
		paramsBytes := []jen.Code{jen.Id("data").Any()}

		file.Func().
			Params(jen.Id("t").Add(jen.Op("*")).Id("templates")).
			Id(methodName).
			Params(params...).
			Error().
			Block(
				jen.Return(
					jen.Id("t").Dot(fieldName).Dot("Execute").Params(
						jen.Id("wr"),
						jen.Id("data"),
					),
				),
			)

		file.Func().
			Params(jen.Id("t").Add(jen.Op("*")).Id("templates")).
			Id(methodBytesName).
			Params(paramsBytes...).
			Params(
				jen.Index().Byte(),
				jen.Error(),
			).
			Block(
				jen.Id("wr").Op(":=").Op("&").Qual("bytes", "Buffer").Values(),
				jen.Id("err").Op(":=").Id("t").Dot(methodName).Params(
					jen.Id("wr"),
					jen.Id("data"),
				),
				jen.Return(jen.List(
					jen.Id("wr").Dot("Bytes").Params(),
					jen.Id("err"),
				)),
			)
	}
}

func (g generator) generateStructConstructor(file *jen.File, templates map[string]Template) error {
	todo := make([]string, 0)
	templatesIter := MapOrdered(templates)

	for name, template := range templatesIter {
		// This is the chain of the template, its base (if any), the base of its base
		// (if any), and so on.
		chain := []string{name}

		for template.Base != "" {
			base, ok := templates[template.Base]
			if !ok {
				return BaseNotFoundError{
					base:     template.Base,
					template: name,
				}
			}

			if slices.Contains(chain, template.Base) {
				return BaseLoopError{
					base:     template.Base,
					template: chain[0],
				}
			}
			chain = append(chain, template.Base)
			template = base
		}

		todo = append(todo, chain...)
	}

	done := make(map[string]struct{})

	file.Func().Id("NewTemplates").Params().Id("Templates").BlockFunc(func(jg *jen.Group) {
		for _, name := range slices.Backward(todo) {
			if _, ok := done[name]; ok {
				continue
			}

			done[name] = struct{}{}
			tmplVarName := fmt.Sprintf("%sTmpl", g.getName(name))
			template := templates[name]

			if template.Base != "" {
				baseTmplVarName := fmt.Sprintf("%sTmpl", g.getName(template.Base))
				jg.Id(tmplVarName).Op(":=").Qual("html/template", "Must").Params(
					jen.Qual("html/template", "Must").Params(
						jen.Id(baseTmplVarName).Dot("Clone").Params(),
					).Dot("Parse").Params(jen.Lit(template.Content)),
				)
			} else {
				jg.Id(tmplVarName).Op(":=").Qual("html/template", "Must").Params(
					jen.Qual("html/template", "New").Params(
						jen.Lit(name),
					).Dot("Parse").Params(
						jen.Lit(template.Content),
					),
				)
			}
		}

		jg.Return(jen.Op("&").Id("templates").Values(jen.DictFunc(func(d jen.Dict) {
			for name := range templates {
				name := g.getName(name)
				fieldName := fmt.Sprintf("template%s", g.caser.String(name))
				tmplVarName := fmt.Sprintf("%sTmpl", g.getName(name))
				d[jen.Id(fieldName)] = jen.Id(tmplVarName)
			}
		})))
	})

	return nil
}

// getName returns the name of the template, without the suffix ".html" if present.
func (g generator) getName(templateName string) string {
	if len(templateName) > 5 && templateName[len(templateName)-5:] == ".html" {
		templateName = templateName[:len(templateName)-5]
	}

	return templateName
}

// DirectoryWriter is an interface to create file in a directy.
type DirectoryWriter interface {
	Create(filename string) (io.WriteCloser, error)
}

// OsDirectoryWriter is an implementation of DirectoryWriter using the "os" library.
type OsDirectoryWriter struct {
	directory string
}

// NewOsDirectoryWriter returns an OsDirectoryWirter object.
//
// It will write to the given directory, which means any path given to its
// Create method will be prefixed by the directory path.
//
// It does not check for directory escape. A path containing ".." can writer in
// another location than the given directory.
func NewOsDirectoryWriter(directory string) (DirectoryWriter, error) {
	generator := OsDirectoryWriter{directory: directory}
	err := generator.createDirectory()
	if err != nil {
		return nil, err
	}

	return generator, nil
}

// Create creates a new file.
//
// It returns an io.WriterCloser object. It is the responsability of the caller to call
// Close() on it.
func (odw OsDirectoryWriter) Create(filename string) (io.WriteCloser, error) {
	file, err := os.Create(path.Join(odw.directory, filename))
	if err != nil {
		return nil, fmt.Errorf("error while creating file %s: %w", filename, err)
	}

	return file, nil
}

func (odw OsDirectoryWriter) createDirectory() error {
	return os.MkdirAll(odw.directory, 0o755)
}

// MapOrdered returns a iterator over the map, ordered by keys.
//
// It assumes no keys are added not deleted to the map after the iterator is created.
// If that's the case, the added keys will be missed, and the deleted keys will be yield.
func MapOrdered(templates map[string]Template) iter.Seq2[string, Template] {
	// we know the templates map will never change, so we order the keys only once
	keys := slices.Sorted(maps.Keys(templates))

	return func(yield func(string, Template) bool) {
		for _, key := range keys {
			if !yield(key, templates[key]) {
				return
			}
		}
	}
}
