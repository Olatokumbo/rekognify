package main

import (
	"rekognify/infra/queue"
	"rekognify/infra/storage"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		name := "rekognify"

		bucket, err := storage.CreateBucket(ctx, name)

		if err != nil {
			return err
		}

		queue, err := queue.CreateQueue(ctx, name, bucket)

		if err != nil {
			return err
		}

		err = storage.CreateS3Notification(ctx, name, bucket, queue)

		if err != nil {
			return err
		}

		storage.ExportS3Outputs(ctx, bucket)

		return nil
	})
}
