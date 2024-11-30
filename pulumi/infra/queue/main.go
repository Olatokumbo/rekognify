package queue

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateQueue(ctx *pulumi.Context, name string, bucket *s3.BucketV2) (*sqs.Queue, error) {

	sqsQueue, err := sqs.NewQueue(ctx, fmt.Sprintf("%s-sqs", name), &sqs.QueueArgs{
		VisibilityTimeoutSeconds: pulumi.Int(60),
	})

	if err != nil {
		return nil, err
	}

	queuePolicyDoc := iam.GetPolicyDocumentOutput(ctx, iam.GetPolicyDocumentOutputArgs{
		Statements: iam.GetPolicyDocumentStatementArray{
			&iam.GetPolicyDocumentStatementArgs{
				Effect: pulumi.String("Allow"),
				Principals: iam.GetPolicyDocumentStatementPrincipalArray{
					&iam.GetPolicyDocumentStatementPrincipalArgs{
						Type: pulumi.String("*"),
						Identifiers: pulumi.StringArray{
							pulumi.String("*"),
						},
					},
				},
				Actions: pulumi.StringArray{
					pulumi.String("sqs:SendMessage"),
				},
				Resources: pulumi.StringArray{
					sqsQueue.Arn,
				},
				Conditions: iam.GetPolicyDocumentStatementConditionArray{
					&iam.GetPolicyDocumentStatementConditionArgs{
						Test:     pulumi.String("ArnEquals"),
						Variable: pulumi.String("aws:SourceArn"),
						Values: pulumi.StringArray{
							bucket.Arn,
						},
					},
				},
			},
		},
	}, nil)

	_, err = sqs.NewQueuePolicy(ctx, name+"-sqs-policy", &sqs.QueuePolicyArgs{
		QueueUrl: sqsQueue.ID(),
		Policy:   queuePolicyDoc.Json(),
	}, pulumi.DependsOn([]pulumi.Resource{sqsQueue}))

	if err != nil {
		return nil, err
	}

	return sqsQueue, nil
}
