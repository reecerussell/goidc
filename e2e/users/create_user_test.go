package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	testData := map[string]string{
		"email":    "test@email.com",
		"password": "MyPassword123",
	}
	userId := ""

	t.Cleanup(func() {
		_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(os.Getenv("USERS_TABLE_NAME")),
			Key: map[string]*dynamodb.AttributeValue{
				"userId": {
					S: aws.String(userId),
				},
			},
		})
		if err != nil {
			panic(err)
		}
	})

	c := &http.Client{
		Timeout: time.Second * 10,
	}

	baseUrl := os.Getenv("BASE_API_URL")
	targetUrl := fmt.Sprintf("%s/test/api/users", baseUrl)
	data, _ := json.Marshal(testData)

	req, _ := http.NewRequest(http.MethodPost, targetUrl, bytes.NewBuffer(data))
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	t.Run("Status Code Should Be OK", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Response Should Contain UserId", func(t *testing.T) {
		var data map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		t.Logf("Body: %v\n", data)

		id, ok := data["id"]
		assert.True(t, ok)

		userId = id.(string)
	})
}

func TestCreateUser_WhereUserAlreadyExists_ReturnsBadRequest(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	testData := map[string]string{
		"userId":   "20340jsnflsfg",
		"email":    "test@email.com",
		"password": "MyPassword123",
	}

	item, _ := dynamodbattribute.MarshalMap(testData)
	_, err := db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("USERS_TABLE_NAME")),
		Item:      item,
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(os.Getenv("USERS_TABLE_NAME")),
			Key: map[string]*dynamodb.AttributeValue{
				"userId": {
					S: aws.String(testData["userId"]),
				},
			},
		})
		if err != nil {
			panic(err)
		}
	})

	c := &http.Client{
		Timeout: time.Second * 10,
	}

	baseUrl := os.Getenv("BASE_API_URL")
	targetUrl := fmt.Sprintf("%s/test/api/users", baseUrl)
	data, _ := json.Marshal(testData)

	req, _ := http.NewRequest(http.MethodPost, targetUrl, bytes.NewBuffer(data))
	req.Header["Content-Type"] = []string{"application/json"}

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	t.Run("Status Code Should Be BadRequest", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Response Should Contain Error", func(t *testing.T) {
		var data map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&data)
		assert.NoError(t, err)

		t.Logf("Body: %v\n", data)

		_, ok := data["error"]
		assert.True(t, ok)
	})
}
