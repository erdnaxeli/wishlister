# statictemplates

`statictemplates` is a tool to generate static go code from templates.

## Why

When using go templates, you have to parse some templates (from either files, FS
or a plain string) at the runtime. This can cause errors if the templates definitions
are incorrect.

Then you need to store those templates somewhere in order to acces them later when
needed. You can do something like this:

```
package statictemplates
import "html/template"
func main() {
  templates := make(map[string]*template.Template)
  // Read all templates. This can be improving by iterating over all file in the
  // directory.
  templates["index.html"] = template.Must(template.ParseFiles("templates/index.html"))
  templates["user.html"] = template.Must(template.ParseFiles("templates/user.html"))
  // later somewhere, render a template
  _ := templates["user.html"].Execute(wr, data)
}
```

This has the disavantage than when executing the template, you refer the template
a string. If I wrote "user.htm" instead, my program would have crashed.

Instead, you may want define some struct with a field per template. The struct would
be filled on the startup of your app, and you would refer to a template using for
example `templates.UserTmpl`. If there were any typo, it would crash on compile time.
This solution is what I prefer, but it is far from perfect. First, you now have to add
a line in the struct and in the method filling it each time you add a new template
(you cannot iterate over the directory files anymore). Second, we still can have an
error when parsing a template.

## What

This package statictemplates answer to those two issues. It generates the struct
definition, and the code to fill it. It also validates that the templates are valid
and copies their content into the generated code to ensure that the parsing step will
not fail.

The only remaining source of errors is the execution step.

Actually, statictemplates generates an interfaces, so you can replace it with a fake
implementation on your test. The interface provides two method for each template
found:

* `Render<TemplateName>(wr io.Writer, data any) error`
* `Render<TemplateName>Bytes(data any) ([]byte, error)`

`<TemplateName>` is the name of the template file, capitalized (the first letter is
changed to its upper version). If the name of the template file ends with ".html",
this suffix is removed.

## Usage

```
go run github.com/erdnaxeli/statictemplates directory/to/templates/ packagename output/directory/
```

## Todo

This package solves 2 issues: ensuring the parsing of the template will not fail and
calling a template statically (instead of using strings). But executing the template
still use the `any` type and can return errors.

A solution would be to infer the types of the data structures to which are applied the
templates, generate the corresponding types, and use them in the methods signatures.
