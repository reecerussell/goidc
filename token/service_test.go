package token

import (
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/reecerussell/gojwt"
	"github.com/reecerussell/gojwt/mock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/goidc/util"
)

func TestGenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlg := mock.NewMockAlgorithm(ctrl)
	mockAlg.EXPECT().Name().Return("RS256", nil)
	mockAlg.EXPECT().Size().Return(256, nil)

	mockSignature := make([]byte, 256)
	rand.Read(mockSignature)

	mockAlg.EXPECT().Sign(gomock.Any()).Return(mockSignature, nil)

	testIssuer := "test"
	testClaims := map[string]interface{}{
		"foo": "bar",
		"one": 1,
	}
	testAudience := "testing"
	testExpirySeconds := int64(3600)

	svc := New(testIssuer)
	token, err := svc.GenerateToken(mockAlg, testClaims, testExpirySeconds, testAudience)
	assert.NoError(t, err)
	assert.Equal(t, "Bearer", token.TokenType)
	assert.Equal(t, testExpirySeconds, token.Expires)

	jwt, err := gojwt.Token(token.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, testAudience, jwt.Claims["aud"])
	assert.Equal(t, testIssuer, jwt.Claims["iss"])
	assert.Equal(t, "bar", jwt.Claims["foo"])
	assert.Equal(t, float64(1), jwt.Claims["one"])

	assert.Equal(t, float64(util.Time().UnixNano()/1e9), jwt.Claims["iat"])
	assert.Equal(t, float64(util.Time().UnixNano()/1e9), jwt.Claims["nbf"])
	assert.Equal(t, float64(util.Time().Add(time.Duration(testExpirySeconds)*time.Second).UnixNano()/1e9), jwt.Claims["exp"])
}

func TestGenerateToken_WhereTheAlgorithmFails_ReturnsError(t *testing.T) {
	testError := errors.New("test error")
	testIssuer := "test"
	testClaims := map[string]interface{}{
		"foo": "bar",
		"one": 1,
	}
	testAudience := "testing"
	testExpirySeconds := int64(3600)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAlg := mock.NewMockAlgorithm(ctrl)
	mockAlg.EXPECT().Name().Return("RS256", nil)
	mockAlg.EXPECT().Size().Return(256, nil)
	mockAlg.EXPECT().Sign(gomock.Any()).Return(nil, testError)

	svc := New(testIssuer)
	token, err := svc.GenerateToken(mockAlg, testClaims, testExpirySeconds, testAudience)
	assert.Nil(t, token)
	assert.Equal(t, testError, err)
}
