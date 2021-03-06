package validator

import (
	"errors"

	"github.com/reecerussell/goidc/dal"
	"github.com/reecerussell/goidc/util"
)

// Validation errors.
var (
	ErrInvalidSecret      = errors.New("invalid client secret")
	ErrInvalidGrantType   = errors.New("invalid grant type")
	ErrMissingScope       = errors.New("missing scope")
	ErrInvalidScope       = errors.New("invalid scope")
	ErrMissingRedirectUri = errors.New("missing redirect uri")
	ErrInvalidRedirectUri = errors.New("invalid redirect uri")
)

// ClientValidator is used to centralize client validation logic, for
// validating incoming requests.
type ClientValidator interface {
	ValidateTokenRequest(c *dal.Client, secret, grantType string, scopes []string) error
	ValidateLoginRequest(c *dal.Client, redirectUri string, scopes []string) error
}

// clientValidator is an implementation of ClientValidator.
type clientValidator struct{}

// NewClientValidator returns a new instance of ClientValidator.
func NewClientValidator() ClientValidator {
	return &clientValidator{}
}

func (*clientValidator) ValidateTokenRequest(c *dal.Client, secret, grantType string, scopes []string) error {
	err := validateSecret(c.Secrets, secret)
	if err != nil {
		return err
	}

	err = validateGrantTypes(c.GrantTypes, grantType)
	if err != nil {
		return err
	}

	err = validateScopes(c.Scopes, scopes)
	if err != nil {
		return err
	}

	return nil
}

func (*clientValidator) ValidateLoginRequest(c *dal.Client, redirectUri string, scopes []string) error {
	if redirectUri == "" {
		return ErrMissingRedirectUri
	}

	if len(scopes) < 1 {
		return ErrMissingScope
	}

	err := validateRedirectUri(c.RedirectUris, redirectUri)
	if err != nil {
		return err
	}

	err = validateScopes(c.Scopes, scopes)
	if err != nil {
		return err
	}

	return nil
}

func validateSecret(allowedSecrets []string, secret string) error {
	for _, allowed := range allowedSecrets {
		if allowed == util.Sha256(secret) {
			return nil
		}
	}

	return ErrInvalidSecret
}

func validateGrantTypes(allowedTypes []string, grantType string) error {
	for _, allowed := range allowedTypes {
		if allowed == grantType {
			return nil
		}
	}

	return ErrInvalidGrantType
}

// Returns an error of any value in scopes is not contained in allowedScopes.
func validateScopes(allowedScopes []string, scopes []string) error {
	for _, scope := range scopes {
		isAllowed := false

		for _, allowed := range allowedScopes {
			if scope == allowed {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return ErrInvalidScope
		}
	}

	return nil
}

func validateRedirectUri(allowed []string, redirectUri string) error {
	for _, uri := range allowed {
		if uri == redirectUri {
			return nil
		}
	}

	return ErrInvalidRedirectUri
}
