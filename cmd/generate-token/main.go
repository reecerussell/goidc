package main

import (
	"encoding/json"
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
		payload := map[string]string{"error": "method not allowed"}
		bytes, _ := json.Marshal(payload)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			IsBase64Encoded: false,
			Body:            string(bytes),
		}, nil
	}

	if req.Headers["Content-Type"] != "application/x-www-form-urlencoded" {
		payload := map[string]string{"error": "invalid content type"}
		bytes, _ := json.Marshal(payload)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			IsBase64Encoded: false,
			Body:            string(bytes),
		}, nil
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
			payload := map[string]string{"error": "invalid client id"}
			bytes, _ := json.Marshal(payload)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Headers: map[string]string{
					"Content-Type": "application/json; charset=utf-8",
				},
				IsBase64Encoded: false,
				Body:            string(bytes),
			}, nil
		}

		payload := map[string]string{"error": err.Error()}
		bytes, _ := json.Marshal(payload)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			IsBase64Encoded: false,
			Body:            string(bytes),
		}, nil
	}

	err = h.validator.ValidateRequest(client, clientSecret, redirectUri, grantType, scopes)
	if err != nil {
		payload := map[string]string{"error": err.Error()}
		bytes, _ := json.Marshal(payload)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			IsBase64Encoded: false,
			Body:            string(bytes),
		}, nil
	}

	claims := map[string]interface{}{
		"sub":    client.ID,
		"scopes": scopes,
	}

	accessToken, err := h.tokens.GenerateToken(claims, 3600, "goidc")
	if err != nil {
		payload := map[string]string{"error": err.Error()}
		bytes, _ := json.Marshal(payload)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Content-Type": "application/json; charset=utf-8",
			},
			IsBase64Encoded: false,
			Body:            string(bytes),
		}, nil
	}

	bytes, _ := json.Marshal(accessToken)
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		},
		IsBase64Encoded: false,
		Body:            string(bytes),
	}, nil
}
