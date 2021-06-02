package main

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	hasher "github.com/reecerussell/adaptive-password-hasher"

	"github.com/reecerussell/goidc"
	"github.com/reecerussell/goidc/dal"
	"github.com/reecerussell/goidc/dal/dynamo"
	"github.com/reecerussell/goidc/util"
	"github.com/reecerussell/goidc/validator"
)

const (
	// TODO: move this config/logic into a centralized place
	//       so that it can be used elsewhere.
	//
	// NOTE: this is copied from validator/user_validator.go
	iterationCount = 10000
	hashKey        = hasher.HashSHA256
)

func main() {
	log.Println("Starting...")

	sess := session.Must(session.NewSession())
	hsr, _ := hasher.New(iterationCount, hasher.DefaultSaltSize, hasher.DefaultKeySize, hashKey)

	hdlr := &Handler{
		up:  dynamo.NewUserProvider(sess),
		us:  dynamo.NewUserService(sess),
		uv:  validator.NewUserValidator(),
		hsr: hsr,
	}

	lambda.Start(hdlr.Handle)
}

// Handle is a struct which ties together a collection of
// dependencies which are needed to handle a request.
type Handler struct {
	hsr hasher.Hasher
	up  dal.UserProvider
	us  dal.UserService
	uv  validator.UserValidator
}

// RequestModel represents the body of an incoming request.
type RequestModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ResponseModel represents a successful response body.
type ResponseModel struct {
	ID string `json:"id"`
}

// Handle is the handler function used to handle a request.
func (h *Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.HTTPMethod != http.MethodPost {
		log.Printf("Invalid method: %s\n", req.HTTPMethod)
		err := errors.New("method not allowed")
		return util.RespondMethodNotAllowed(err), nil
	}

	if h := util.Header(req, "Content-Type"); strings.Index(h, "application/json") == -1 {
		log.Printf("Invalid content type: %s\n", h)
		err := errors.New("invalid content type")
		return util.RespondBadRequest(err), nil
	}

	var model RequestModel
	util.ReadJSON(req, &model)

	err := h.uv.ValidateUser(model.Email, model.Password)
	if err != nil {
		log.Printf("Invalid user data: %v\n", err)

		return util.RespondBadRequest(err), nil
	}

	ctx = goidc.NewContext(ctx, &req)
	_, err = h.up.GetByEmail(ctx, model.Email)
	if err != dal.ErrUserNotFound {
		if err == nil {
			log.Printf("User already exists: %s\n", model.Email)
			err := errors.New("user already exists")
			return util.RespondBadRequest(err), nil
		}

		log.Printf("users: failed to get user by email: %v\n", err)

		return util.RespondError(err), nil
	}

	passwordHash := h.hsr.Hash([]byte(model.Password))
	user := &dal.User{
		ID:           uuid.New().String(),
		Email:        model.Email,
		PasswordHash: base64.StdEncoding.EncodeToString(passwordHash),
	}
	err = h.us.Create(ctx, user)
	if err != nil {
		log.Printf("users: failed to create user: %v\n", err)

		return util.RespondError(err), nil
	}

	log.Printf("Created user with id: %s\n", user.ID)

	data := ResponseModel{ID: user.ID}
	return util.Respond(http.StatusOK, data), nil
}
