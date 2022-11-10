package uploader

import (
	"context"
	"io/ioutil"
	"testing"
)

//func newBucket() *Uploader {
//	return New("maou-iekai", "AKIA25FPY244AVYYZ2", "jUPM39nASDkch4FxiI8d83tkhYsXhP2oZOdURW", "ap-southeast-1", 2)
//}
//
//func TestUploadFile(t *testing.T) {
//	uploader := newBucket()
//	bytes, err := ioutil.ReadFile("example/rectangular_logo.png")
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	if _, err := uploader.Upload(context.Background(), "banners/rectangular_logo.png", bytes); err != nil {
//		t.Error(err)
//	}
//}
//
//func TestGetURLFile(t *testing.T) {
//	uploader := newBucket()
//	url, err := uploader.GetURL(context.Background(), "banners/47ae5f99-11b0-4864-a016-42660bd22bc5.jpg")
//	if err != nil {
//		t.Error(err)
//	}
//
//	t.Log(url)
//}
//
//func TestDeleteFile(t *testing.T) {
//	uploader := newBucket()
//	if err := uploader.Delete(context.Background(), "banners/rectangular_logo.png"); err != nil {
//		t.Error(err)
//	}
//}

func newBucket() *Uploader {
	config := &Config{
		Bucket:       "zsmartex",
		AccessKey:    "DO00KNY6QY9WY74X3ULM",
		AccessSecret: "1oRKb2zH8H2VcgR4goLnC6r4DNQr7TfeSWsUKUCRQbo",
		Region:       "sgp1",
		Endpoint:     "sgp1.digitaloceanspaces.com",
		Version:      1,
	}
	return New(config)
}

func TestUploadFile(t *testing.T) {
	uploader := newBucket()
	bytes, err := ioutil.ReadFile("example/rectangular_logo.png")
	if err != nil {
		t.Error(err)
		return
	}

	if err := uploader.Upload(context.Background(), "banners/rectangular_logo.png", bytes); err != nil {
		t.Error(err)
	}
}

func TestGetURLFile(t *testing.T) {
	uploader := newBucket()
	url, err := uploader.GetURL(context.Background(), "banners/rectangular_logo.png")
	if err != nil {
		t.Error(err)
	}

	t.Log(url)
}

func TestGetBodyContent(t *testing.T) {
	uploader := newBucket()
	url, err := uploader.GetBodyContent(context.Background(), "banners/rectangular_logo.png")
	if err != nil {
		t.Error(err)
	}

	t.Log(url)
}

func TestDeleteFile(t *testing.T) {
	uploader := newBucket()
	if err := uploader.Delete(context.Background(), "banners/rectangular_logo.png"); err != nil {
		t.Error(err)
	}
}
