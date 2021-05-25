package validator

import (
	"encoding/base64"
	"errors"

	hasher "github.com/reecerussell/adaptive-password-hasher"
	"github.com/reecerussell/goidc/dal"
)

const (
	iterationCount = 10000
	hashKey        = hasher.HashSHA256
)

// Common validation errors.
var (
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidPasswordHash = errors.New("invalid password hash")
)

// UserValidator is used to centralize user validation logic.
type UserValidator interface {
	// ValidatePassword is used to validate a given user's password.
	ValidatePassword(u *dal.User, password string) error
}

type userValidator struct {
	h hasher.Hasher
}

// NewUserValidator returns a new instance of UserValidator.
func NewUserValidator() UserValidator {
	h, _ := hasher.New(iterationCount, hasher.DefaultSaltSize, hasher.DefaultKeySize, hashKey)

	return &userValidator{
		h: h,
	}
}

func (v *userValidator) ValidatePassword(u *dal.User, password string) error {
	if password == "" {
		return ErrInvalidPassword
	}

	bytes, err := base64.StdEncoding.DecodeString(u.PasswordHash)
	if err != nil {
		return ErrInvalidPasswordHash
	}

	if !v.h.Verify([]byte(password), bytes) {
		return ErrInvalidPassword
	}

	return nil
}
