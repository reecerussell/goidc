package dal

// UserService is used to perform write-operations
// on the users domain.
type UserService interface {
	// Create inserts a user record into the data store.
	Create(u *User) error
}
