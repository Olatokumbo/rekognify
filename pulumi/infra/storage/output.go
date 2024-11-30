package storage

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ExportS3Outputs(ctx *pulumi.Context, s3 *s3.BucketV2) {
	ctx.Export("bucket-name", s3.ID())
	ctx.Export("bucket-arn", s3.Arn)
}
