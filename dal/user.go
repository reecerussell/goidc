package dal

// User represents the structure of a user in the database.
type User struct {
	ID           string `json:"userId"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}
