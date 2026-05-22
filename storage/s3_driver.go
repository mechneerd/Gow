package storage

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Driver implements the Driver interface for Amazon S3 using AWS SDK v2.
type S3Driver struct {
	client *s3.Client
	bucket string
	region string
}

// NewS3Driver creates a new S3 driver.
// It supports loading credentials from environment variables or explicit keys.
func NewS3Driver(bucket, region, accessKey, secretKey string) (*S3Driver, error) {
	var cfg aws.Config
	var err error

	ctx := context.Background()

	if accessKey != "" && secretKey != "" {
		// Use explicit credentials
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	} else {
		// Use default credential chain (env, IAM role, etc.)
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
		)
	}

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Driver{
		client: client,
		bucket: bucket,
		region: region,
	}, nil
}

func (d *S3Driver) Put(path string, contents io.Reader) error {
	_, err := d.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
		Body:   contents,
	})
	return err
}

func (d *S3Driver) Get(path string) (io.ReadCloser, error) {
	output, err := d.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (d *S3Driver) Delete(path string) error {
	_, err := d.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	return err
}

func (d *S3Driver) Exists(path string) bool {
	_, err := d.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	})
	return err == nil
}

func (d *S3Driver) URL(path string) string {
	// Return a direct public URL (assumes bucket is public or objects are public)
	// For private objects, you should generate presigned URLs instead.
	return "https://" + d.bucket + ".s3." + d.region + ".amazonaws.com/" + path
}

// PresignedURL generates a temporary signed URL (useful for private files).
func (d *S3Driver) PresignedURL(path string, expiry time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(d.client)

	result, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(path),
	}, s3.WithPresignExpires(expiry))

	if err != nil {
		return "", err
	}
	return result.URL, nil
}
