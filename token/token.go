package token

// Token contains an access token and any relevent data
// about the token, such as, expiry and type.
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int64  `json:"expires"`
}
