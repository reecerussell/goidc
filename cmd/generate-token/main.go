package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/reecerussell/gojwt/kms"

	"github.com/reecerussell/goidc/dal"
	"github.com/reecerussell/goidc/dal/dynamo"
	"github.com/reecerussell/goidc/token"
	"github.com/reecerussell/goidc/util"
	"github.com/reecerussell/goidc/validator"
)

const keyIdVar = "TOKEN_KEY_ID"

func main() {
	log.Println("Starting...")

	keyId := os.Getenv(keyIdVar)
	log.Printf("Key Id: %s\n", keyId)

	sess := session.Must(session.NewSession())
	alg, _ := kms.New(sess, keyId, kms.RSA_PKCS1_S256)

	tokenService := token.New(alg, "goidc")
	clientProvider := dynamo.NewClientProvider(sess)

	hdlr := &Handler{
		tokens:    tokenService,
		clients:   clientProvider,
		validator: validator.NewClientValidator(),
	}

	lambda.Start(hdlr.Handle)
}

type Handler struct {
	tokens    token.Service
	clients   dal.ClientProvider
	validator validator.ClientValidator
}

func (h *Handler) Handle(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != http.MethodPost {
		return util.RespondMethodNotAllowed(errors.New("method not allowed")), nil
	}

	if req.Headers["Content-Type"] != "application/x-www-form-urlencoded" {
		return util.RespondBadRequest(errors.New("invalid content type")), nil
	}

	data, _ := url.ParseQuery(req.Body)
	clientId := data.Get("client_id")
	clientSecret := data.Get("client_secret")
	grantType := data.Get("grant_type")
	redirectUri := data.Get("redirect_uri")
	scopes := strings.Split(" ", data.Get("scope"))

	client, err := h.clients.Get(clientId)
	if err != nil {
		if err == dal.ErrClientNotFound {
			return util.RespondBadRequest(errors.New("invalid client id")), nil
		}

		return util.RespondError(err), nil
	}

	err = h.validator.ValidateRequest(client, clientSecret, redirectUri, grantType, scopes)
	if err != nil {
		return util.RespondBadRequest(err), nil
	}

	claims := map[string]interface{}{
		"sub":    client.ID,
		"scopes": scopes,
	}

	accessToken, err := h.tokens.GenerateToken(claims, 3600, "goidc")
	if err != nil {
		return util.RespondError(err), nil
	}

	return util.RespondOk(accessToken), nil
}
