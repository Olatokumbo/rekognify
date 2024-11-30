package storage

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateS3Notification(ctx *pulumi.Context, name string, bucket *s3.BucketV2, queue *sqs.Queue) error {
	_, err := s3.NewBucketNotification(ctx, fmt.Sprintf("%s-bucket-notification", name), &s3.BucketNotificationArgs{
		Bucket: bucket.ID(),
		Queues: s3.BucketNotificationQueueArray{
			&s3.BucketNotificationQueueArgs{
				QueueArn: queue.Arn,
				Events: pulumi.StringArray{
					pulumi.String("s3:ObjectCreated:Post"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
