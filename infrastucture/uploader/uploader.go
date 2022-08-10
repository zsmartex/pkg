package uploader

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	credentialsV2 "github.com/aws/aws-sdk-go-v2/credentials"
	managerV2 "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	s3V2 "github.com/aws/aws-sdk-go-v2/service/s3"

	awsV1 "github.com/aws/aws-sdk-go/aws"
	credentialsV1 "github.com/aws/aws-sdk-go/aws/credentials"
	sessionV1 "github.com/aws/aws-sdk-go/aws/session"
	s3V1 "github.com/aws/aws-sdk-go/service/s3"
)

type Uploader struct {
	bucket          string
	clientV2        *s3V2.Client
	uploaderV2      *managerV2.Uploader
	presignClientV2 *s3V2.PresignClient

	clientV1   *s3V1.S3
	uploaderV1 *s3V1.S3

	version int64
}

type Config struct {
	Bucket       string
	AccessKey    string
	AccessSecret string
	Region       string
	Enpoint      string
	Version      int64
}

func New(config *Config) *Uploader {
	if config.Version == 2 {
		awsConfig := awsV2.Config{
			Credentials: credentialsV2.NewStaticCredentialsProvider(config.AccessKey, config.AccessSecret, ""),
			Region:      config.Region,
		}

		client := s3V2.NewFromConfig(awsConfig)

		return &Uploader{
			clientV2:        client,
			bucket:          config.Bucket,
			presignClientV2: s3V2.NewPresignClient(client),
			uploaderV2: managerV2.NewUploader(client, func(u *managerV2.Uploader) {
				u.PartSize = 10 * 1024 * 1024 // 10MB per part
				u.Concurrency = 5
			}),
			version: 2,
		}
	}

	// version 1
	s3Config := &awsV1.Config{
		Credentials: credentialsV1.NewStaticCredentials(config.AccessKey, config.AccessSecret, ""),
		Endpoint:    awsV1.String(config.Enpoint),
		Region:      awsV1.String(config.Region),
	}

	newSession := sessionV1.Must(sessionV1.NewSession(s3Config))
	s3Client := s3V1.New(newSession)

	return &Uploader{
		clientV1:   s3Client,
		bucket:     config.Bucket,
		uploaderV1: s3Client,
		version:    1,
	}
}

func (u *Uploader) GetURL(context context.Context, key string) (string, error) {
	if u.version == 2 {
		resp, err := u.presignClientV2.PresignGetObject(context, &s3V2.GetObjectInput{
			Bucket:                     awsV2.String(u.bucket),
			Key:                        awsV2.String(key),
			ResponseContentDisposition: awsV2.String("inline"),
		})

		if err != nil {
			return "", err
		}

		return resp.URL, nil
	}

	// version 1
	req, _ := u.clientV1.GetObjectRequest(&s3V1.GetObjectInput{
		Bucket:                     awsV1.String(u.bucket),
		Key:                        awsV1.String(key),
		ResponseContentDisposition: awsV1.String("inline"),
	})

	urlStr, err := req.Presign(10 * time.Minute)
	if err != nil {
		fmt.Println(err.Error())
	}

	return urlStr, nil
}

func (u *Uploader) GetBodyContent(context context.Context, key string) ([]byte, error) {
	if u.version == 2 {
		resp, err := u.clientV2.GetObject(context, &s3V2.GetObjectInput{
			Bucket: awsV2.String(u.bucket),
			Key:    awsV2.String(key),
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

	// version 1
	input := &s3V1.GetObjectInput{
		Bucket: awsV1.String(u.bucket),
		Key:    awsV1.String(key),
	}

	result, err := u.clientV1.GetObject(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *Uploader) Upload(context context.Context, key string, body []byte) error {
	contentType := http.DetectContentType(body)

	if u.version == 2 {
		_, err := u.uploaderV2.Upload(context, &s3V2.PutObjectInput{
			Bucket:      awsV2.String(u.bucket),
			Key:         awsV2.String(key),
			Body:        bytes.NewReader(body),
			ContentType: awsV2.String(contentType),
		})

		return err
	}

	// version 1
	_, err := u.uploaderV1.PutObjectWithContext(context, &s3V1.PutObjectInput{
		Bucket:      awsV1.String(u.bucket),
		Key:         awsV1.String(key),
		Body:        bytes.NewReader(body),
		ContentType: awsV1.String(contentType),
	})

	return err
}

func (u *Uploader) Delete(context context.Context, key string) error {
	if u.version == 2 {
		_, err := u.clientV2.DeleteObject(context, &s3V2.DeleteObjectInput{
			Bucket: awsV2.String(u.bucket),
			Key:    awsV2.String(key),
		})

		return err
	}

	// version 1
	input := &s3V1.DeleteObjectInput{
		Bucket: awsV1.String(u.bucket),
		Key:    awsV1.String(key),
	}

	_, err := u.clientV1.DeleteObject(input)

	return err
}
