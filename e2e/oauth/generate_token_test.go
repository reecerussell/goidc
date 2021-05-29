package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	db := dynamodb.New(sess)

	testClientId := "9238ulfdsfre"
	ha := sha256.New()
	ha.Write([]byte("32y4i423"))
	testClientSecret := base64.StdEncoding.EncodeToString(ha.Sum(nil))
	testData := map[string]interface{}{
		"clientId":   testClientId,
		"name":       "TestGenerateToken",
		"scopes":     []string{"test"},
		"grantTypes": []string{"client_credentials"},
		"secrets":    []string{testClientSecret},
	}

	av, err := dynamodbattribute.MarshalMap(testData)
	if err != nil {
		panic(err)
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("CLIENTS_TABLE_NAME")),
		Item:      av,
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, err := db.DeleteItem(&dynamodb.DeleteItemInput{
			TableName: aws.String(os.Getenv("CLIENTS_TABLE_NAME")),
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

	// request the api
	c := &http.Client{
		Timeout: time.Second * 10,
	}

	baseUrl := os.Getenv("BASE_API_URL")
	targetUrl := fmt.Sprintf("%s/test/oauth/token", baseUrl)

	reqData := url.Values{
		"client_id":     {testClientId},
		"client_secret": {"32y4i423"},
		"grant_type":    {"client_credentials"},
		"scope":         {"test"},
	}
	body := strings.NewReader(reqData.Encode())

	req, err := http.NewRequest(http.MethodPost, targetUrl, body)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(body.Len()))

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	t.Run("Status Code Should Be OK", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Response Should Contain Token", func(t *testing.T) {
		var tokenData map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&tokenData)
		assert.NoError(t, err)

		t.Logf("Body: %v\n", tokenData)

		assert.Equal(t, "Bearer", tokenData["token_type"])
		assert.Equal(t, float64(3600), tokenData["expires"])

		token := tokenData["access_token"].(string)
		payloadB64 := strings.Split(token, ".")[1]
		payloadBytes, _ := base64.RawURLEncoding.DecodeString(payloadB64)

		var payload map[string]interface{}
		err = json.Unmarshal(payloadBytes, &payload)

		assert.Equal(t, []string{"test"}, payload["scopes"])
	})
}
