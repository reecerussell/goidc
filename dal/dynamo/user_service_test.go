package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/goidc/dal"
)

func TestCreateUser(t *testing.T) {
	ctx := buildUsersContext()
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	testUser := &dal.User{
		ID:           "239lqjrpw3e",
		Email:        "john@doe.com",
		PasswordHash: "38974enflndf",
	}

	t.Cleanup(func() {
		_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(UsersTableName(ctx)),
			Key: map[string]*dynamodb.AttributeValue{
				"userId": {
					S: aws.String(testUser.ID),
				},
			},
		})
		if err != nil {
			panic(err)
		}
	})

	s := NewUserService(sess)
	err := s.Create(ctx, testUser)
	assert.NoError(t, err)

	t.Run("User Should Be Created", func(t *testing.T) {
		res, err := db.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(UsersTableName(ctx)),
			Key: map[string]*dynamodb.AttributeValue{
				"userId": {
					S: aws.String(testUser.ID),
				},
			},
		})
		if err != nil {
			panic(err)
		}

		var user dal.User
		err = dynamodbattribute.UnmarshalMap(res.Item, &user)
		if err != nil {
			panic(err)
		}

		assert.Equal(t, testUser.ID, user.ID)
		assert.Equal(t, testUser.Email, user.Email)
		assert.Equal(t, testUser.PasswordHash, user.PasswordHash)
	})
}
