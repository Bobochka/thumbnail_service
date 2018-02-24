package store

import (
	"bytes"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	bucket     *string
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
}

func New(endpoint, region, bucket string) (*S3, error) {
	cfg := &aws.Config{Region: aws.String(region)}
	if endpoint != "" {
		cfg.Endpoint = aws.String(endpoint)
	}

	s, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	svc := s3.New(s)

	return &S3{
		bucket:     aws.String(bucket),
		uploader:   s3manager.NewUploaderWithClient(svc),
		downloader: s3manager.NewDownloaderWithClient(svc),
	}, nil
}

func (s *S3) Get(key string) []byte {
	params := &s3.GetObjectInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
	}

	res := []byte{}
	buffer := aws.NewWriteAtBuffer(res)
	_, err := s.downloader.Download(buffer, params)

	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == s3.ErrCodeNoSuchKey {
			return buffer.Bytes()
		}
	}

	if err != nil {
		log.Printf("unable to read from s3: %+v\n", err)
	}

	return buffer.Bytes()
}

func (s *S3) Set(key string, data []byte) error {
	buf := bytes.NewReader(data)

	params := &s3manager.UploadInput{
		Bucket: s.bucket,
		Key:    aws.String(key),
		Body:   buf,
	}

	_, err := s.uploader.Upload(params)

	return err
}
