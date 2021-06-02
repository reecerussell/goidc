package validator

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	hasher "github.com/reecerussell/adaptive-password-hasher"

	"github.com/reecerussell/goidc/dal"
)

const (
	// TODO: move this config/logic into a centralized place
	//       so that it can be used elsewhere.
	iterationCount = 10000
	hashKey        = hasher.HashSHA256

	minPasswordLength = 6
	emailPattern      = "[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,6}"
)

// Common validation errors.
var (
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrEmailEmpty          = errors.New("email cannot be empty")
	ErrEmailInvalid        = errors.New("email is not valid")
	ErrPasswordEmpty       = errors.New("password cannot be empty")
	ErrPasswordTooShort    = fmt.Errorf("password must be at least %d characters long", minPasswordLength)
)

// UserValidator is used to centralize user validation logic.
type UserValidator interface {
	// ValidatePassword is used to validate a given user's password.
	ValidatePassword(u *dal.User, password string) error

	// ValidateUser is used to validate core user values,
	// such as email, password etc. Should be used when
	// creating users.
	ValidateUser(email, password string) error
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

func (v *userValidator) ValidateUser(email, password string) error {
	if email == "" {
		return ErrEmailEmpty
	}

	re := regexp.MustCompile(emailPattern)
	if !re.MatchString(email) {
		return ErrEmailInvalid
	}

	if password == "" {
		return ErrPasswordEmpty
	}

	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	return nil
}
