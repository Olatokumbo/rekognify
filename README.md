# Rekognify

Rekognify is an image recognition tool developed using AWS SAM and Golang.

## Architecture

The Rekognify architecture consists of the following components:

- **AWS S3**: For storing uploaded images.
- **AWS Lambda**: Processes images and generates recognition results.
- **AWS SQS**: Buffer zone for incoming tasks.
- **AWS API Gateway**: Provides endpoints for uploading images and retrieving results.
- **AWS Rekognition**: Performs image classification and labeling.
- **AWS Cloudfront**: Content Delivery Network to accelerate media delivery.

## API Usage Instructions

### 1. Upload an Image to the S3 Bucket via a Signed URL

You can use the following cURL command to generate an S3 PreSigned URL and upload the image directly to the S3 bucket:

**Request:**

```bash
curl --location 'https://api.rekognify.com/upload' \
--header 'Content-Type: application/json' \
--data '{
    "filename": "test.jpg",
    "mimeType": "image/jpeg"
}'
```

**Response:**

```json
{
    "url": "https://presigned-s3-url",
    "id": "test-uuid.jpg"
}
```

- `url`: The S3 PreSigned URL.
- `id`: The unique name or ID of the uploaded image.

**Supported MIME Types:**

- `image/jpeg`
- `image/png`
- `image/gif`
- `image/webp`

**Note:** Ensure the uploaded image matches the specified MIME type.

**Upload the Image:**

Use the `PUT` method with the PreSigned URL to upload the image. A `200 OK` status should be returned upon success.

### 2. Retrieve Classification Results

Query the API to retrieve the classification labels for the uploaded image.

**Request:**

```bash
curl --location 'https://api.rekognify.com/info/{test-uuid}.jpg'
```

**Response:**

```json
{
    "filename": "{test-uuid}.jpg",
    "url": "https://cdn.rekognify.com/{test-uuid}.jpg",
    "labels": [
        {
            "category": "Weapons and Military",
            "confidence": 97.92,
            "name": "Launch"
        },
        {
            "category": "Weapons and Military",
            "confidence": 86.16,
            "name": "Weapon"
        }
    ]
}
```

- `filename`: The name of the uploaded image.
- `url`: The CDN-hosted URL for the uploaded image.
- `labels`: An array of classification results with `category`, `confidence`, and `name` for each identified label.


## Development

* AWS CLI already configured with Administrator permission.
* [Docker installed](https://www.docker.com/community-edition).
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html).
* [Golang](https://golang.org).

### Setup Process

#### Installing Dependencies and Building the Target

The built-in `sam build` command is used to build a Docker image from a Dockerfile and then copy the source of your application into the Docker image.
