package dal

// Client represents the structure of a client in the database.
type Client struct {
	ID           string   `json:"clientId"`
	Name         string   `json:"name"`
	RedirectUris []string `json:"redirectUris"`
	GrantTypes   []string `json:"grantTypes"`
	Scopes       []string `json:"scopes"`
	Secrets      []string `json:"secrets"`
}
