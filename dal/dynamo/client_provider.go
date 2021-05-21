package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/reecerussell/goidc/dal"
)

// ClientProvider is an implementation of dal.ClientProvider fo DynamoDB.
type ClientProvider struct {
	svc *dynamodb.DynamoDB
}

// NewClientProvider returns a new instance of ClientProvider,
// for the given session, sess.
func NewClientProvider(sess *session.Session) dal.ClientProvider {
	return &ClientProvider{
		svc: dynamodb.New(sess),
	}
}

// Get queries the clients table in DynamoDB for a client with the given id.
func (p *ClientProvider) Get(ctx context.Context, id string) (*dal.Client, error) {
	res, err := p.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(ClientsTableName(ctx)),
		Key: map[string]*dynamodb.AttributeValue{
			"clientId": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if res.Item == nil {
		return nil, dal.ErrClientNotFound
	}

	var client dal.Client
	err = dynamodbattribute.UnmarshalMap(res.Item, &client)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
