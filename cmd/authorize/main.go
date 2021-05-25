package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/reecerussell/gojwt"
	"github.com/reecerussell/gojwt/kms"

	"github.com/reecerussell/goidc"
	"github.com/reecerussell/goidc/dal"
	"github.com/reecerussell/goidc/dal/dynamo"
	"github.com/reecerussell/goidc/token"
	"github.com/reecerussell/goidc/util"
	"github.com/reecerussell/goidc/validator"
)

var (
	errInvalidCredentials      = errors.New("email and/or password is invalid")
	errUnsupportedResponseType = errors.New("unsupported response type")
)

func main() {
	log.Println("Starting...")

	sess := session.Must(session.NewSession())

	hdlr := &Handler{
		sess:      sess,
		tokens:    token.New("goidc"),
		clients:   dynamo.NewClientProvider(sess),
		clientVal: validator.NewClientValidator(),
		users:     dynamo.NewUserProvider(sess),
		userVal:   validator.NewUserValidator(),
	}

	lambda.Start(hdlr.Handle)
}

// Handler is used to provide a Lambda handler function.
type Handler struct {
	sess      *session.Session
	tokens    token.Service
	users     dal.UserProvider
	userVal   validator.UserValidator
	clients   dal.ClientProvider
	clientVal validator.ClientValidator
}

// LoginModel represents the body of the login request.
type LoginModel struct {
	ClientID     string   `json:"clientId"`
	RedirectUri  string   `json:"redirectUri"`
	Scopes       []string `json:"scopes"`
	ResponseType string   `json:"responseType"`
	State        string   `json:"state"`
	Nonce        string   `json:"nonce"`

	Email    string `json:"email"`
	Password string `json:"password"`
}

// ResponseModel represents a successfull request's response body.
type ResponseModel struct {
	RedirectUri string `json:"redirectUri"`
}

func (h *Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != http.MethodPost {
		err := errors.New("method not allowed")
		return util.RespondMethodNotAllowed(err), nil
	}

	if strings.Index(req.Headers["Content-Type"], "application/json") == -1 {
		err := errors.New("invalid content type")
		return util.RespondBadRequest(err), nil
	}

	var model LoginModel
	err := json.Unmarshal([]byte(req.Body), &model)
	if err != nil {
		err = errors.New("invalid content")
		return util.RespondBadRequest(err), nil
	}

	ctx = goidc.NewContext(ctx, &req)
	client, err := h.clients.Get(ctx, model.ClientID)
	if err != nil {
		if err == dal.ErrClientNotFound {
			err = errors.New("invalid client id")
			return util.RespondBadRequest(err), nil
		}

		return util.RespondError(err), nil
	}

	err = h.clientVal.ValidateLoginRequest(client, model.RedirectUri, model.Scopes)
	if err != nil {
		return util.RespondBadRequest(err), nil
	}

	user, err := h.users.GetByEmail(ctx, model.Email)
	if err != nil {
		if err == dal.ErrUserNotFound {
			return util.RespondBadRequest(errInvalidCredentials), nil
		}

		return util.RespondError(err), nil
	}

	err = h.userVal.ValidatePassword(user, model.Password)
	if err != nil {
		if err == validator.ErrInvalidPassword {
			return util.RespondBadRequest(errInvalidCredentials), nil
		}

		return util.RespondError(err), nil
	}

	switch model.ResponseType {
	case "id_token token":
		return h.idTokenTokenResponse(ctx, client, user, &model)
	default:
		return util.RespondBadRequest(errUnsupportedResponseType), nil
	}
}

func (h *Handler) idTokenTokenResponse(ctx context.Context, c *dal.Client, u *dal.User, m *LoginModel) (events.APIGatewayProxyResponse, error) {
	alg, _ := kms.New(h.sess, goidc.StageVariable(ctx, "JWT_KEY_ID"), kms.RSA_PKCS1_S256)
	jwt, err := h.generateAccessToken(alg, u.Email)
	if err != nil {
		return util.RespondError(err), nil
	}

	idToken, err := h.generateIdToken(alg, u.ID, m.State, &jwt.AccessToken)
	if err != nil {
		return util.RespondError(err), nil
	}

	urlValues := url.Values{
		"id_token":     {idToken},
		"state":        {m.State},
		"nonce":        {m.Nonce},
		"access_token": {jwt.AccessToken},
		"token_type":   {jwt.TokenType},
		"expires_in":   {strconv.Itoa(int(jwt.Expires))},
	}

	redirectUri := fmt.Sprintf("%s?%s", m.RedirectUri, urlValues.Encode())
	resp := ResponseModel{RedirectUri: redirectUri}

	return util.Respond(http.StatusOK, resp), nil
}

func (h *Handler) generateIdToken(alg gojwt.Algorithm, sub, state string, accessToken *string) (string, error) {
	claims := map[string]interface{}{
		"sub":     sub,
		"s_hash":  util.Sha256Half(state),
		"at_hash": util.Sha256Half(*accessToken),
	}

	jwt, err := h.tokens.GenerateToken(alg, claims, 36000, "goidc")
	if err != nil {
		return "", err
	}

	return jwt.AccessToken, nil
}

func (h *Handler) generateAccessToken(alg gojwt.Algorithm, sub string) (*token.Token, error) {
	claims := map[string]interface{}{
		"sub": sub,
	}

	return h.tokens.GenerateToken(alg, claims, 3600, "goidc")
}
