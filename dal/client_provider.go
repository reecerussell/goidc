package dal

import "errors"

// ErrClientNotFound is a common error used when a
// client cannot be found, or does not exist.
var ErrClientNotFound = errors.New("client not found")

// ClientProvider is used to retrieve client information from the database.
type ClientProvider interface {
	// Get retrieves a client from the database, with the given id.
	// If the client cannot be found, ErrClientNotFound will be returned
	// as the error.
	Get(id string) (*Client, error)
}
