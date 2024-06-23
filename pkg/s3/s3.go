package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Config struct {
	Region   string
	Endpoint string
	Bucket   string
	ID       string
	Secret   string
	Timeout  time.Duration
}

type S3Bucket struct {
	config S3Config
	s3     *s3.S3
}

func NewWithConfig(config S3Config) *S3Bucket {

	// All clients require a Session. The Session provides the client with
	// shared configuration such as region, endpoint, and credentials. A
	// Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess := session.Must(session.NewSession(
		aws.NewConfig().
			WithRegion(config.Bucket).
			WithEndpoint(config.Endpoint).
			WithCredentials(
				credentials.NewStaticCredentials(config.ID, config.Secret, ""),
			).WithS3ForcePathStyle(true),
	))

	// Create a new instance of the service's client with a Session.
	// Optional aws.Config values can also be provided as variadic arguments
	// to the New function. This option allows you to provide service
	// specific configuration.
	svc := s3.New(sess)

	return &S3Bucket{config, svc}
}

func (s *S3Bucket) Upload(ctx context.Context, path string, body io.ReadSeeker) error {
	_, err := s.s3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(path),
		Body:   body,
	})

	if err != nil {
		return err
	}

	fmt.Printf("successfully uploaded file to %s/%s\n", s.config.Bucket, path)
	return nil
}

func (s *S3Bucket) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := s.s3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(path),
	})
	return output.Body, err
}
