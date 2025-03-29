package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type RequestPayload struct {
	Mimetype string `json:"mimetype"`
}

type ResponsePayload struct {
	URL string `json:"url"`
	ID  string `json:"id"`
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "OPTIONS,POST,GET,PUT",
	"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
	"Content-Type":                 "application/json",
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var payload RequestPayload

	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to parse request body: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
			Headers:    corsHeaders,
		}, nil
	}

	if payload.Mimetype == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Mimetype is required.",
			StatusCode: http.StatusBadRequest,
			Headers:    corsHeaders,
		}, nil
	}

	if !allowedMimeTypes[payload.Mimetype] {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Unsupported file type: %s", payload.Mimetype),
			StatusCode: http.StatusBadRequest,
			Headers:    corsHeaders,
		}, nil
	}

	bucket := os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		return events.APIGatewayProxyResponse{
			Body:       "S3_BUCKET_NAME environment variable not set.",
			StatusCode: http.StatusInternalServerError,
			Headers:    corsHeaders,
		}, nil
	}

	bucketPrefix := os.Getenv("S3_BUCKET_PREFIX")
	if bucketPrefix == "" {
		return events.APIGatewayProxyResponse{
			Body:       "S3_BUCKET_PREFIX environment variable not set.",
			StatusCode: http.StatusInternalServerError,
			Headers:    corsHeaders,
		}, nil
	}

	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to create AWS session: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
			Headers:    corsHeaders,
		}, nil
	}

	s3SVC := s3.New(sess)

	exts, err := mime.ExtensionsByType(payload.Mimetype)
	if err != nil || len(exts) == 0 {
		fmt.Println("Unknown mimetype, defaulting to .jpg")
		exts = []string{".jpg"}
	}
	ext := strings.TrimPrefix(exts[0], ".")

	filename := fmt.Sprintf("%s.%s", uuid.New().String(), ext)

	req, _ := s3SVC.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", bucketPrefix, filename)),
		ContentType: aws.String(payload.Mimetype),
	})

	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to generate presigned URL: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
			Headers:    corsHeaders,
		}, nil
	}

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if tableName == "" {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("DYNAMODB_TABLE_NAME variable not set: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
			Headers:    corsHeaders,
		}, nil
	}

	dynamodbSess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynamodbSVC := dynamodb.New(dynamodbSess)

	_, err = dynamodbSVC.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]*dynamodb.AttributeValue{
			"filename": {
				S: aws.String(filename),
			},
			"status": {
				S: aws.String("READY_TO_BE_PROCESSED"),
			},
			"created_at": {
				S: aws.String(time.Now().Format(time.RFC3339)),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error inserting item into DynamoDB: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
			Headers:    corsHeaders,
		}, nil
	}

	response := ResponsePayload{
		URL: url,
		ID:  filename,
	}
	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: http.StatusOK,
		Headers:    corsHeaders,
	}, nil
}

func main() {
	lambda.Start(handler)
}
