package types

const UserContextKey = "user"

type AuthenticatedUser struct {
	ID       string
	Email    string
	LoggedIn bool
}
