package api

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (a *Api) signedUploadUrl(key string, svc *s3.S3) string {
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_UPLOAD_BUCKET")),
		Key:    aws.String(key),
	})

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		a.Logger.Error().Err(err).Msg("failed to sign request")
	}

	return url
}

func (a *Api) s3Client() *s3.S3 {
	sess, err := session.NewSession(a.awsConfig())
	if err != nil {
		a.Logger.Error().Err(err).Msg("unable to create a new S3 session")
	}

	svc := s3.New(sess)

	return svc
}

func (a *Api) awsConfig() *aws.Config {
	return &aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	}
}
