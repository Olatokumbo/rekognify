package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type S3Event struct {
	Records []struct {
		S3 struct {
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func handler(ctx context.Context, event events.SQSEvent) (events.SQSEventResponse, error) {
	var failures []events.SQSBatchItemFailure

	for _, record := range event.Records {
		fmt.Printf("Processing Message: %s\n", record.MessageId)
		fmt.Printf("SQS Message Body: %s\n", record.Body)

		var s3Event S3Event
		err := json.Unmarshal([]byte(record.Body), &s3Event)
		if err != nil {
			fmt.Printf("Error unmarshalling SQS message: %v\n", err)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		if len(s3Event.Records) == 0 {
			fmt.Println("No records found in S3 event")
			continue
		}

		filename := s3Event.Records[0].S3.Object.Key

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := dynamodb.New(sess)

		tableName := os.Getenv("DYNAMODB_TABLE_NAME")

		if tableName == "" {
			fmt.Printf("DYNAMODB_TABLE_NAME variable not set: %v\n", err)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		_, err = svc.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]*dynamodb.AttributeValue{
				"filename": {
					S: aws.String(filename),
				},
				"status": {
					S: aws.String("PENDING"),
				},
				"created_at": {
					S: aws.String(time.Now().Format(time.RFC3339)),
				},
			},
		})
		if err != nil {
			fmt.Printf("Error inserting item into DynamoDB: %v\n", err)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
		}
	}

	return events.SQSEventResponse{
		BatchItemFailures: failures,
	}, nil
}

func main() {
	lambda.Start(handler)
}
