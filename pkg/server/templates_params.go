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

// UserListsViewList represents a wishlist in the UserListsView template.
type UserListsViewList struct {
	AdminID string
	ID      string
	Name    string
}

// ParamsUserListsView holds the parameters for the UserListsView template.
type ParamsUserListsView struct {
	Lists []UserListsViewList

	Error string
}
