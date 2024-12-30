package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ResponsePayload struct {
	Filename string  `json:"filename"`
	Labels   []Label `json:"labels"`
}

type Label struct {
	Category   string  `json:"category"`
	Confidence float64 `json:"confidence"`
	Name       string  `json:"name"`
}

func parseLabels(labels []*dynamodb.AttributeValue) []Label {
	var parsedLabels []Label
	for _, label := range labels {
		fmt.Println(label)
		parsedLabels = append(parsedLabels, Label{
			Category:   *label.M["category"].S,
			Confidence: parseFloat(*label.M["confidence"].N),
			Name:       *label.M["name"].S,
		})
	}
	return parsedLabels
}

func parseFloat(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return f
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	filename := request.PathParameters["filename"]
	if filename == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Filename is required",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	dynamo := dynamodb.New(sess)

	imageData, err := dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"filename": {
				S: aws.String(filename),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Error getting item from DynamoDB",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if imageData.Item == nil {
		return events.APIGatewayProxyResponse{
			Body:       "Item not found",
			StatusCode: http.StatusNotFound,
		}, nil
	}

	if *imageData.Item["status"].S != "COMPLETED" {
		return events.APIGatewayProxyResponse{
			Body:       "Item is still being processed",
			StatusCode: http.StatusConflict,
		}, nil
	}

	response := ResponsePayload{
		Filename: filename,
		Labels:   parseLabels(imageData.Item["labels"].L),
	}

	responseBody, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}