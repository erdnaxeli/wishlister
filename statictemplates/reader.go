package statictemplates

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"text/template/parse"
)

// Template represent a template.
type Template struct {
	// Base is the name of base template of the template. If the template does not have
	// a base, it default to "".
	Base string
	// Content is the full content of the template.
	Content string
}

// TemplateReader is the main interface to read multiples templates.
//
// Where the files are read from depends on the implementation.
type TemplateReader interface {
	// ReadTemplates read templates and return a map of Template object.
	//
	// The keys of the map are the names of the templates.
	ReadTemplates() (map[string]Template, error)
}

type templateReader struct {
	dir fs.FS
}

// NewTemplateReader returns a TemplateReader that reads templates from a directory.
//
// It reads all files in the given directory. If the directory contains subdirectories,
// they will be ignored.
func NewTemplateReader(directory string) TemplateReader {
	dir := os.DirFS(directory)
	return NewTemplateReaderFS(dir)
}

// NewTemplateReaderFS returns a TemplateReader that reads templates from an fs.FS object.
//
// It reads all files in the given FS. If the FS contains subdirectories, they will be
// ignored.
func NewTemplateReaderFS(dir fs.FS) TemplateReader {
	return templateReader{
		dir: dir,
	}
}

func (tr templateReader) ReadTemplates() (map[string]Template, error) {
	templates := make(map[string]Template)

	err := fs.WalkDir(tr.dir, ".", func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			if err != nil {
				// It means there was an error when calling Stat() on the root dir.
				return fmt.Errorf("error while reading templates directory: %w", err)
			}

			return nil
		}

		if d.IsDir() {
			// We ignore any directories (for now).
			return fs.SkipDir
		}

		content, err := tr.readFile(path)
		if err != nil {
			return fmt.Errorf("error while reading template %s: %w", d.Name(), err)
		}

		baseTemplateName, err := tr.readBaseTemplateName(content)
		if err != nil {
			return fmt.Errorf("error while reading template %s: %w", d.Name(), err)
		}

		templates[d.Name()] = Template{
			Base:    baseTemplateName,
			Content: content,
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func (tr templateReader) readFile(path string) (string, error) {
	content, err := fs.ReadFile(tr.dir, path)
	if err != nil {
		return "", fmt.Errorf("error while reading file %s: %w", path, err)
	}

	return string(content), nil
}

// readBaseTemplateNmae read the base template name (if any) of the current template.
//
// The base template name must be defined in a template comment like this:
// ```
// {{/* base: someTemplate */}}
// ```
//
// The comment must be the first line in the template.
func (tr templateReader) readBaseTemplateName(content string) (string, error) {
	// We parse the template with text/template/parse because text/template does not
	// provide a way to keep comments in the nodes tree.
	// html/template uses text/template behind the scene.
	// We also need to not check that functions are defined (SkipFuncCheck) because
	// those functions are added by the text/template module and not exported.
	tree := parse.Tree{Mode: parse.ParseComments | parse.SkipFuncCheck}
	treeSet := make(map[string]*parse.Tree)
	_, err := tree.Parse(content, "", "", treeSet)
	if err != nil {
		return "", fmt.Errorf("error while parsing template: %w", err)
	}

	if len(tree.Root.Nodes) == 0 {
		// The template is empty, wich is not very useful, but that's not a reason
		// to crash. The user probably knows what they are doing.
		println("no nodes")
		return "", nil
	}

	node := tree.Root.Nodes[0]
	if node.Type() != parse.NodeComment {
		return "", nil
	}

	comment := node.(*parse.CommentNode).Text
	if comment[:9] == "/* base: " {
		name := comment[9 : len(comment)-2]
		return strings.TrimSpace(name), nil
	}

	return "", nil
}
