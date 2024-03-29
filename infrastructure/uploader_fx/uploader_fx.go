package uploader_fx

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
	"github.com/zsmartex/pkg/v2/config"
	"go.uber.org/fx"

	awsV1 "github.com/aws/aws-sdk-go/aws"
	credentialsV1 "github.com/aws/aws-sdk-go/aws/credentials"
	sessionV1 "github.com/aws/aws-sdk-go/aws/session"
	s3V1 "github.com/aws/aws-sdk-go/service/s3"
)

var (
	Module = fx.Module(
		"uploader_fx",
		fx.Provide(
			New,
		),
	)
)

type Client struct {
	bucket          string
	clientV2        *s3V2.Client
	uploaderV2      *managerV2.Uploader
	presignClientV2 *s3V2.PresignClient

	clientV1   *s3V1.S3
	uploaderV1 *s3V1.S3

	version int64
}

type uploaderParams struct {
	fx.In

	Config config.ObjectStorage
}

func New(params uploaderParams) *Client {
	if params.Config.Version == 2 {
		awsConfig := awsV2.Config{
			Credentials: credentialsV2.NewStaticCredentialsProvider(params.Config.AccessKey, params.Config.AccessSecret, ""),
			Region:      params.Config.Region,
		}

		client := s3V2.NewFromConfig(awsConfig)

		return &Client{
			clientV2:        client,
			bucket:          params.Config.Bucket,
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
		Credentials: credentialsV1.NewStaticCredentials(params.Config.AccessKey, params.Config.AccessSecret, ""),
		Endpoint:    awsV1.String(params.Config.Endpoint),
		Region:      awsV1.String(params.Config.Region),
	}

	newSession := sessionV1.Must(sessionV1.NewSession(s3Config))
	s3Client := s3V1.New(newSession)

	return &Client{
		clientV1:   s3Client,
		bucket:     params.Config.Bucket,
		uploaderV1: s3Client,
		version:    1,
	}
}

func (u *Client) GetURL(context context.Context, key string) (string, error) {
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

func (u *Client) GetBodyContent(context context.Context, key string) ([]byte, error) {
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

func (u *Client) Upload(context context.Context, key string, body []byte) error {
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

func (u *Client) Delete(context context.Context, key string) error {
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
