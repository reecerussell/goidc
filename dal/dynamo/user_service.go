package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/reecerussell/goidc/dal"
)

// UserService is an implementation of dal.UserService for DynamoDB.
type UserService struct {
	svc *dynamodb.DynamoDB
}

// NewUserService returns a new instance of UserService.
func NewUserService(sess *session.Session) dal.UserService {
	return &UserService{
		svc: dynamodb.New(sess),
	}
}

func (s *UserService) Create(ctx context.Context, u *dal.User) error {
	item, _ := dynamodbattribute.MarshalMap(u)

	_, err := s.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(UsersTableName(ctx)),
		Item:      item,
	})
	if err != nil {
		return err
	}

	return nil
}
