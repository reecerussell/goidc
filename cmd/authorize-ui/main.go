package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/reecerussell/goidc"
	"github.com/reecerussell/goidc/util"
)

func main() {

}

type Handler struct {
	sess *session.Session
}

func (h *Handler) Handle(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	downloader := s3manager.NewDownloader(h.sess)

	path := strings.Replace(req.Path, "/oauth/authorize", "", 1)
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if path == "" {
		path = "index.html"
	}

	ctx = goidc.NewContext(ctx, &req)
	version := goidc.StageVariable(ctx, "UI_VERSION")
	path = fmt.Sprintf("%s/%s", version, path)

	buf := aws.NewWriteAtBuffer([]byte{})
	n, err := downloader.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(goidc.StageVariable(ctx, "UI_BUCKET")),
		Key:    aws.String(path),
	})
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
			}, nil
		}

		return util.RespondError(err), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		IsBase64Encoded: true,
		Headers: map[string]string{
			"Content-Length": strconv.Itoa(int(n)),
			"Content-Type":   http.DetectContentType(buf.Bytes()),
		},
		Body: base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}
