package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/goidc/dal"
	"github.com/reecerussell/goidc/util"
)

func TestClientValidator_ValidateTokenRequest_ReturnsNoError(t *testing.T) {
	testClient := &dal.Client{
		RedirectUris: []string{"http://localhost:8080"},
		GrantTypes:   []string{"client_credentials"},
		Scopes:       []string{"openid", "email", "profile"},
		Secrets:      []string{"239473204", util.Sha256("test")},
	}

	cv := NewClientValidator()
	err := cv.ValidateTokenRequest(testClient, "test", "client_credentials", []string{"openid"})
	assert.NoError(t, err)
}

func TestClientValidator_ValidateTokenRequest_ReturnsError(t *testing.T) {
	testClient := &dal.Client{
		RedirectUris: []string{"http://localhost:8080"},
		GrantTypes:   []string{"client_credentials"},
		Scopes:       []string{"openid", "email", "profile"},
		Secrets:      []string{"239473204", util.Sha256("test")},
	}

	cv := NewClientValidator()

	t.Run("Given Invalid Secret", func(t *testing.T) {
		err := cv.ValidateTokenRequest(testClient, "hello", "client_credentials", []string{"openid"})
		assert.Equal(t, ErrInvalidSecret, err)
	})

	t.Run("Given Invalid GrantType", func(t *testing.T) {
		err := cv.ValidateTokenRequest(testClient, "test", "code", []string{"openid"})
		assert.Equal(t, ErrInvalidGrantType, err)
	})

	t.Run("Given Invalid Scope", func(t *testing.T) {
		err := cv.ValidateTokenRequest(testClient, "test", "client_credentials", []string{"openid", "test"})
		assert.NotNil(t, err)
	})
}
