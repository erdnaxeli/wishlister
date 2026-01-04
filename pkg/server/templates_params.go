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

// ParamsLogin holds the parameters for the Login template.
type ParamsLogin struct {
	Email string

	Error      string
	EmailError string

	Sent bool
}
