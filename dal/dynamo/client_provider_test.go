package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	sess := session.Must(session.NewSession())
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
		TableName: aws.String(ClientsTableName()),
		Item:      av,
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(ClientsTableName()),
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
		client, err := cp.Get(testClientId)
		assert.NoError(t, err)
		assert.Equal(t, testClientId, client.ID)
		assert.Equal(t, testData["name"], client.Name)
		assert.Contains(t, client.RedirectUris, testData["redirectUris"])
		assert.Equal(t, 1, len(client.RedirectUris))
		assert.Contains(t, client.Scopes, testData["scopes"])
		assert.Equal(t, 1, len(client.Scopes))
	})
}
