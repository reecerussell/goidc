package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

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
	filter := expression.Name("email").Equal(expression.Value(email))
	projection := expression.NamesList(expression.Name("userId"), expression.Name("email"), expression.Name("passwordHash"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return nil, err
	}

	res, err := p.svc.Scan(&dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(UsersTableName()),
	})
	if err != nil {
		return nil, err
	}

	if len(res.Items) < 1 {
		return nil, dal.ErrUserNotFound
	}

	var user dal.User
	err = dynamodbattribute.UnmarshalMap(res.Items[0], &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
