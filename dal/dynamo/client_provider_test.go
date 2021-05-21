package dynamo

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/goidc"
)

func buildClientsContext() context.Context {
	req := events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"CLIENTS_TABLE_NAME": "goidc-clients-test",
		},
	}

	return goidc.NewContext(context.Background(), &req)
}

func TestGetClient(t *testing.T) {
	ctx := buildClientsContext()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	testClientId := "9238ulfdsfre"
	testData := map[string]interface{}{
		"clientId":     testClientId,
		"name":         "TestGetClient",
		"redirectUris": []string{"http://localhost:3000"},
		"scopes":       []string{"test"},
	}

	av, err := dynamodbattribute.MarshalMap(testData)
	if err != nil {
		panic(err)
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(ClientsTableName(ctx)),
		Item:      av,
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(ClientsTableName(ctx)),
			Key: map[string]*dynamodb.AttributeValue{
				"clientId": {
					S: aws.String(testClientId),
				},
			},
		})
		if err != nil {
			panic(err)
		}
	})

	t.Run("Client Should Be Returned", func(t *testing.T) {
		cp := NewClientProvider(sess)
		client, err := cp.Get(ctx, testClientId)
		assert.NoError(t, err)
		assert.Equal(t, testClientId, client.ID)
		assert.Equal(t, testData["name"], client.Name)
		assert.Equal(t, client.RedirectUris, testData["redirectUris"])
		assert.Equal(t, client.Scopes, testData["scopes"])
	})
}
