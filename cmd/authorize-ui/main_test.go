package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/reecerussell/goidc/util"
	"github.com/stretchr/testify/assert"
)

func TestHandle_WhereFileExists_ReturnsFileAsBase64(t *testing.T) {
	testFileData := []byte(`<!DOCTYPE html><html><head><title>Hello World</title></head><body><h1>Hello World</h1></body></html>`)
	testVersion := "integration-tests"
	testFileName := "index.html"
	s3Key := fmt.Sprintf("%s/%s", testVersion, testFileName)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	uploader := s3manager.NewUploader(sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("UI_BUCKET")),
		Key:    aws.String(s3Key),
		Body:   bytes.NewBuffer(testFileData),
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		svc := s3.New(sess)
		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("UI_BUCKET")),
			Key:    aws.String(s3Key),
		})
		if err != nil {
			panic(err)
		}

		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(os.Getenv("UI_BUCKET")),
			Key:    aws.String(s3Key),
		})
		if err != nil {
			panic(err)
		}
	})

	h := &Handler{sess: sess}
	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		Path: "/oauth/authorize/",
		StageVariables: map[string]string{
			"UI_BUCKET":  os.Getenv("UI_BUCKET"),
			"UI_VERSION": testVersion,
		},
	})
	assert.NoError(t, err)

	t.Run("Should Return OK", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Should Return Content Type", func(t *testing.T) {
		assert.Contains(t, resp.Headers["Content-Type"], "text/html")
		assert.Equal(t, strconv.Itoa(len(testFileData)), resp.Headers["Content-Length"])
		assert.Equal(t, "public, max-age=604800", resp.Headers["Cache-Control"])
		assert.Equal(t, util.Sha256(string(testFileData)), resp.Headers["ETag"])
	})

	t.Run("Should Return Content", func(t *testing.T) {
		assert.True(t, resp.IsBase64Encoded)

		data, err := base64.StdEncoding.DecodeString(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(testFileData), string(data))
	})
}

func TestHandle_WhereFileDoesNotExist_ReturnsNotFound(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	h := &Handler{sess: sess}
	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		Path: "/oauth/authorize/notafile.txt",
		StageVariables: map[string]string{
			"UI_BUCKET":  os.Getenv("UI_BUCKET"),
			"UI_VERSION": "not-a-version",
		},
	})
	assert.NoError(t, err)

	t.Run("Should Return NotFound", func(t *testing.T) {
		t.Logf("Response: %v\n", resp.Body)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestHandle_WhereS3ReturnsError_ReturnsInternalServerError(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	h := &Handler{sess: sess}
	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		Path: "/oauth/authorize/notafile.txt",
		StageVariables: map[string]string{
			"UI_BUCKET":  "not-a-bucket",
			"UI_VERSION": "not-a-version",
		},
	})
	assert.NoError(t, err)

	t.Run("Should Return InternalServerError", func(t *testing.T) {
		t.Logf("Response: %v\n", resp.Body)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Should Return Error", func(t *testing.T) {
		assert.Equal(t, "application/json; charset=utf-8", resp.Headers["Content-Type"])
		assert.False(t, resp.IsBase64Encoded)

		var data map[string]interface{}
		err := json.Unmarshal([]byte(resp.Body), &data)
		assert.NoError(t, err)

		_, ok := data["error"]
		assert.True(t, ok)
	})
}

func TestHandle_GivenUndetectableContentType_ReturnsFileAsBase64(t *testing.T) {
	testFileData := []byte(`html,body{padding:0;}`)
	testVersion := "integration-tests"
	testFileName := "styles.css"
	s3Key := fmt.Sprintf("%s/%s", testVersion, testFileName)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	uploader := s3manager.NewUploader(sess)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("UI_BUCKET")),
		Key:    aws.String(s3Key),
		Body:   bytes.NewBuffer(testFileData),
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		svc := s3.New(sess)
		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(os.Getenv("UI_BUCKET")),
			Key:    aws.String(s3Key),
		})
		if err != nil {
			panic(err)
		}

		err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
			Bucket: aws.String(os.Getenv("UI_BUCKET")),
			Key:    aws.String(s3Key),
		})
		if err != nil {
			panic(err)
		}
	})

	h := &Handler{sess: sess}
	ctx := context.Background()
	resp, err := h.Handle(ctx, events.APIGatewayProxyRequest{
		Path: "/oauth/authorize/" + testFileName,
		StageVariables: map[string]string{
			"UI_BUCKET":  os.Getenv("UI_BUCKET"),
			"UI_VERSION": testVersion,
		},
	})
	assert.NoError(t, err)

	t.Run("Should Return OK", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Should Return Content Type", func(t *testing.T) {
		assert.Contains(t, resp.Headers["Content-Type"], "text/css")
		assert.Equal(t, strconv.Itoa(len(testFileData)), resp.Headers["Content-Length"])
	})

	t.Run("Should Return Content", func(t *testing.T) {
		assert.True(t, resp.IsBase64Encoded)

		data, err := base64.StdEncoding.DecodeString(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, string(testFileData), string(data))
	})
}

func TestGetFileTypeByExtension(t *testing.T) {
	testData := map[string]string{
		"test.js":        "application/javascript; charset=utf-8",
		"styles.min.css": "text/css; charset=utf-8",
		".my-file-type":  "text/plain; charset=utf-8",
		"index.html":     "text/html; charset=utf-8",
	}

	for file, expected := range testData {
		assert.Equal(t, expected, getFileTypeByExtension(file))
	}
}
