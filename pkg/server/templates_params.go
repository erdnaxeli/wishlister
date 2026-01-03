package server

// ParamsNew holds the parameters for the New template.
type ParamsNew struct {
	Name  string
	User  string
	Email string

	Error      string
	NameError  string
	UserError  string
	EmailError string
}
