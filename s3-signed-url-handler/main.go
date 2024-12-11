package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type RequestPayload struct {
	Filename string `json:"filename"`
	Mimetype string `json:"mimetype"`
}

type ResponsePayload struct {
	URL string `json:"url"`
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var payload RequestPayload

	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to parse request body: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if payload.Filename == "" || payload.Mimetype == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Both 'filename' and 'mimetype' are required.",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if !allowedMimeTypes[payload.Mimetype] {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Unsupported file type: %s", payload.Mimetype),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		return events.APIGatewayProxyResponse{
			Body:       "S3_BUCKET environment variable not set.",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to create AWS session: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fmt.Sprintf("%s_%s", uuid.New().String(), payload.Filename)),
		ContentType: aws.String(payload.Mimetype),
	})

	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to generate presigned URL: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	response := ResponsePayload{
		URL: url,
	}
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
