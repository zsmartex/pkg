package uploader

import (
	"io/ioutil"
	"testing"
)

func newBucket() *Uploader {
	return New("zsmartex-tech", "AKIA6MAKWE2NUXH5KAG4", "bUSTAgwMFbjME1BNR/bpWy4l0IPX+TlhXpqNOc+Q", "us-east-1")
}

func TestUploadFile(t *testing.T) {
	uploader := newBucket()
	bytes, err := ioutil.ReadFile("example/rectangular_logo.png")
	if err != nil {
		t.Error(err)
		return
	}

	uploader.Upload("banners/rectangular_logo.png", bytes)
}

func TestGetURLFile(t *testing.T) {
	uploader := newBucket()
	uploader.GetURL("banners/rectangular_logo.png")
}
