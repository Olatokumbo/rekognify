package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/rekognition"
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

		tableName := os.Getenv("DYNAMODB_TABLE_NAME")

		if tableName == "" {
			fmt.Printf("DYNAMODB_TABLE_NAME variable not set: %v\n", err)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		bucketName := os.Getenv("S3_BUCKET")
		if bucketName == "" {
			fmt.Printf("S3_BUCKET variable not set: %v\n", err)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		if tableName == "" {
			fmt.Printf("DYNAMODB_TABLE_NAME variable not set: %v\n", err)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
			continue
		}

		filename := s3Event.Records[0].S3.Object.Key

		dynamodbSess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		dynamodbSVC := dynamodb.New(dynamodbSess)

		_, err = dynamodbSVC.PutItem(&dynamodb.PutItemInput{
			TableName: &tableName,
			Item: map[string]*dynamodb.AttributeValue{
				"filename": {
					S: aws.String(filename),
				},
				"status": {
					S: aws.String("PROCESSING"),
				},
			},
		})

		if err != nil {
			fmt.Printf("Error updating item from DynamoDB: %v\n", err)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
		}
		rekognitionSess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		rekognitionSVC := rekognition.New(rekognitionSess)

		labelResult, err := rekognitionSVC.DetectLabels(&rekognition.DetectLabelsInput{
			Image: &rekognition.Image{
				S3Object: &rekognition.S3Object{
					Bucket: aws.String(bucketName),
					Name:   &filename,
				},
			},
			MaxLabels: aws.Int64(10),
		})

		if err != nil {
			fmt.Printf("Error detecting labels: %v\n", err)
			failures = append(failures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
		}

		saveLabelsToDynamoDB(filename, labelResult.Labels, tableName, dynamodbSVC)

	}

	return events.SQSEventResponse{
		BatchItemFailures: failures,
	}, nil
}

func main() {
	lambda.Start(handler)
}

func saveLabelsToDynamoDB(fileName string, labelResult []*rekognition.Label, tableName string, dynamodbSVC *dynamodb.DynamoDB) error {
	var labels []*dynamodb.AttributeValue

	for _, label := range labelResult {
		labelItem := map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(*label.Name),
			},
			"category": {
				S: aws.String(*label.Categories[0].Name),
			},
			"confidence": {
				N: aws.String(fmt.Sprintf("%.2f", *label.Confidence)),
			},
		}

		labels = append(labels, &dynamodb.AttributeValue{M: labelItem})
	}

	item := map[string]*dynamodb.AttributeValue{
		"filename": {
			S: aws.String(fileName),
		},
		"status": {
			S: aws.String("COMPLETED"),
		},
		"labels": {
			L: labels,
		},
	}

	_, err := dynamodbSVC.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item into DynamoDB: %w", err)
	}
	return nil
}
