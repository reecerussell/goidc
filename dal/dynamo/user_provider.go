package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/reecerussell/goidc/dal"
)

// UserProvider is an implementation of dal.UserProvider for DynamoDB.
type UserProvider struct {
	svc *dynamodb.DynamoDB
}

// NewUserProvider returns a new instance of UserProvider for the given session, sess.
func NewUserProvider(sess *session.Session) dal.UserProvider {
	return &UserProvider{
		svc: dynamodb.New(sess),
	}
}

// GetByEmail queries the users DynamoDB table for a user with the given email.
func (p *UserProvider) GetByEmail(email string) (*dal.User, error) {
	res, err := p.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(UsersTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if res.Item == nil {
		return nil, dal.ErrUserNotFound
	}

	var user dal.User
	err = dynamodbattribute.UnmarshalMap(res.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
