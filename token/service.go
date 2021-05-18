package token

import (
	"time"

	"github.com/reecerussell/goidc/util"

	"github.com/reecerussell/gojwt"
)

// Service is a high level interface used to generate and
// verify JSON-Web Tokens.
type Service interface {
	GenerateToken(claims map[string]interface{}, expirySeconds int64, audience string) (*Token, error)
}

type service struct {
	alg    gojwt.Algorithm
	issuer string
}

func New(alg gojwt.Algorithm, tokenIssuer string) Service {
	return &service{
		alg:    alg,
		issuer: tokenIssuer,
	}
}

func (s *service) GenerateToken(claims map[string]interface{}, expirySeconds int64, audience string) (*Token, error) {
	now := util.Time()
	expiry := now.Add(time.Duration(expirySeconds) * time.Second)

	builder, _ := gojwt.New(s.alg)
	jwt, err := builder.AddClaims(claims).
		AddClaim("iss", s.issuer).
		AddClaim("aud", audience).
		SetExpiry(expiry).
		SetIssuedAt(now).
		SetNotBefore(now).
		Build()
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken: jwt,
		TokenType:   "Bearer",
		Expires:     expirySeconds,
	}, nil
}
