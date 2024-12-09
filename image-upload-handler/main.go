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

// RequestPayload defines the expected JSON structure
type RequestPayload struct {
	Filename string `json:"filename"`
	Mimetype string `json:"mimetype"`
}

// ResponsePayload defines the structure for the response
type ResponsePayload struct {
	URL string `json:"url"`
}

// Allowed MIME types for images
var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the JSON body
	var payload RequestPayload
	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to parse request body: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	// Validate required fields
	if payload.Filename == "" || payload.Mimetype == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Both 'filename' and 'mimetype' are required.",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	// Validate mimetype
	if !allowedMimeTypes[payload.Mimetype] {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Unsupported file type: %s", payload.Mimetype),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	// Get S3 bucket name from environment variables
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		return events.APIGatewayProxyResponse{
			Body:       "S3_BUCKET environment variable not set.",
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	// Create an AWS session
	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to create AWS session: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	// Create an S3 service client
	svc := s3.New(sess)

	// Generate a presigned URL
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fmt.Sprintf("%s_%s", uuid.New().String(), payload.Filename)),
		ContentType: aws.String(payload.Mimetype),
	})

	url, err := req.Presign(5 * time.Minute) // URL valid for 5 minutes
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to generate presigned URL: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	// Return the presigned URL
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
