package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/reecerussell/goidc"
	"github.com/reecerussell/gojwt/kms"

	"github.com/reecerussell/goidc/dal"
	"github.com/reecerussell/goidc/dal/dynamo"
	"github.com/reecerussell/goidc/token"
	"github.com/reecerussell/goidc/util"
	"github.com/reecerussell/goidc/validator"
)

func main() {
	log.Println("Starting...")

	sess := session.Must(session.NewSession())
	tokenService := token.New("goidc")
	clientProvider := dynamo.NewClientProvider(sess)

	hdlr := &Handler{
		sess:      sess,
		tokens:    tokenService,
		clients:   clientProvider,
		validator: validator.NewClientValidator(),
	}

	lambda.Start(hdlr.Handle)
}

type Handler struct {
	sess      *session.Session
	tokens    token.Service
	clients   dal.ClientProvider
	validator validator.ClientValidator
}

func (h *Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != http.MethodPost {
		return util.RespondMethodNotAllowed(errors.New("method not allowed")), nil
	}

	if util.Header(req, "Content-Type") != "application/x-www-form-urlencoded" {
		log.Printf("Invalid Content Type: %v", req.Headers["Content-Type"])
		return util.RespondBadRequest(errors.New("invalid content type")), nil
	}

	data := util.ReadForm(req)
	clientId := data.Get("client_id")
	clientSecret := data.Get("client_secret")
	grantType := data.Get("grant_type")
	scopes := strings.Split(data.Get("scope"), " ")

	ctx = goidc.NewContext(ctx, &req)
	client, err := h.clients.Get(ctx, clientId)
	if err != nil {
		if err == dal.ErrClientNotFound {
			return util.RespondBadRequest(errors.New("invalid client id")), nil
		}

		return util.RespondError(err), nil
	}

	err = h.validator.ValidateTokenRequest(client, clientSecret, grantType, scopes)
	if err != nil {
		return util.RespondBadRequest(err), nil
	}

	claims := map[string]interface{}{
		"sub":    client.ID,
		"scopes": scopes,
	}

	alg, _ := kms.New(h.sess, goidc.StageVariable(ctx, "JWT_KEY_ID"), kms.RSA_PKCS1_S256)
	accessToken, err := h.tokens.GenerateToken(alg, claims, 3600, "goidc")
	if err != nil {
		return util.RespondError(err), nil
	}

	return util.RespondOk(accessToken), nil
}
