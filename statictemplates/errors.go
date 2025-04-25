package statictemplates

import "fmt"

// BaseNotFoundError is the error when the base of a template is not found.
type BaseNotFoundError struct {
	base     string
	template string
}

func (err BaseNotFoundError) Error() string {
	return fmt.Sprintf("base %s for template %s not found", err.base, err.template)
}

// BaseLoopError is the error when all the bases of a template make a loop.
type BaseLoopError struct {
	base     string
	template string
}

func (err BaseLoopError) Error() string {
	return fmt.Sprintf("loop detected with the base %s of the template %s", err.base, err.template)
}
