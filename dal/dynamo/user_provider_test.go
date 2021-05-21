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

func buildUsersContext() context.Context {
	req := events.APIGatewayProxyRequest{
		StageVariables: map[string]string{
			"USERS_TABLE_NAME": "goidc-users-test",
		},
	}

	return goidc.NewContext(context.Background(), &req)
}

func TestGetUserByEmail(t *testing.T) {
	ctx := buildUsersContext()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	testUserId := "wlerhewrlw"
	testData := map[string]interface{}{
		"userId":       testUserId,
		"email":        "test@test.go",
		"passwordHash": "3wirwhc8o",
	}

	av, err := dynamodbattribute.MarshalMap(testData)
	if err != nil {
		panic(err)
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(UsersTableName(ctx)),
		Item:      av,
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(UsersTableName(ctx)),
			Key: map[string]*dynamodb.AttributeValue{
				"userId": {
					S: aws.String(testUserId),
				},
			},
		})
		if err != nil {
			panic(err)
		}
	})

	t.Run("User Should Be Returned", func(t *testing.T) {
		p := NewUserProvider(sess)
		user, err := p.GetByEmail(ctx, "test@test.go")
		assert.NoError(t, err)
		assert.Equal(t, testUserId, user.ID)
		assert.Equal(t, testData["email"], user.Email)
		assert.Equal(t, testData["passwordHash"], user.PasswordHash)
	})
}
