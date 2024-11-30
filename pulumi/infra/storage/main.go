package storage

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateBucket(ctx *pulumi.Context, name string) (*s3.BucketV2, error) {
	bucket, err := s3.NewBucketV2(ctx, fmt.Sprintf("%s-image-upload-bucket", name), nil)

	if err != nil {
		return nil, err
	}

	return bucket, nil
}
