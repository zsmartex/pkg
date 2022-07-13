package uploader

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Uploader struct {
	bucket        string
	client        *s3.Client
	uploader      *manager.Uploader
	presignClient *s3.PresignClient
}

func New(bucket, accessKey, accessSecret string, region string) *Uploader {
	awsConfig := aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, accessSecret, ""),
		Region:      region,
	}

	client := s3.NewFromConfig(awsConfig)

	return &Uploader{
		client:        client,
		bucket:        bucket,
		presignClient: s3.NewPresignClient(client),
		uploader: manager.NewUploader(client, func(u *manager.Uploader) {
			u.PartSize = 10 * 1024 * 1024 // 10MB per part
			u.Concurrency = 5
		}),
	}
}

func (u *Uploader) GetURL(context context.Context, key string) (string, error) {
	resp, err := u.presignClient.PresignGetObject(context, &s3.GetObjectInput{
		Bucket:                     aws.String(u.bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String("inline"),
	})

	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

func (u *Uploader) GetBodyContent(context context.Context, key string) ([]byte, error) {
	resp, err := u.client.GetObject(context, &s3.GetObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *Uploader) Upload(context context.Context, key string, body []byte) (*manager.UploadOutput, error) {
	content_type := http.DetectContentType(body)

	return u.uploader.Upload(context, &s3.PutObjectInput{
		Bucket:      aws.String(u.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(content_type),
	})
}

func (u *Uploader) Delete(context context.Context, key string) error {
	_, err := u.client.DeleteObject(context, &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
	})

	return err
}
