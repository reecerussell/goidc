package dal

import "errors"

// ErrUserNotFound is a common error used when a user cannot be found, or does not exist.
var ErrUserNotFound = errors.New("user not found")

// UserProvider is a DAL interface used to retrieve user data from the database.
type UserProvider interface {
	GetByEmail(email string) (*User, error)
}
