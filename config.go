package main

import (
	"fmt"
	"os"

	"github.com/Bobochka/thumbnail_service/lib"
	"github.com/Bobochka/thumbnail_service/lib/downloader"
	"github.com/Bobochka/thumbnail_service/lib/locker"
	"github.com/Bobochka/thumbnail_service/lib/service"
	"github.com/Bobochka/thumbnail_service/lib/store"
)

const (
	defaultAwsRegion = "us-east-1"
	defaultBucket    = "cldnrthumbnails"
	defaultRedisURL  = "redis://localhost:6379"
	defaultBindPort  = "8080"
)

func ReadConfig() (*service.Config, error) {
	store, err := store.New(s3Endpoint(), awsRegion(), bucketName())
	if err != nil {
		return nil, fmt.Errorf("unable to init s3 store: %s", err)
	}

	locker, err := locker.New(redisURL())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %s", err)
	}

	return &service.Config{
		Store:      store,
		Downloader: downloader.New(lib.SupportedContentTypes),
		Locker:     locker,
	}, nil
}

func bucketName() string {
	name := os.Getenv("S3_BUCKET_NAME")
	if name == "" {
		name = defaultBucket
	}
	return name
}

func redisURL() string {
	url := os.Getenv("REDIS_URL")
	if url == "" {
		url = defaultRedisURL
	}
	return url
}

func bindPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultBindPort
	}
	return ":" + port
}

func awsRegion() string {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = defaultAwsRegion
	}
	return region
}

func s3Endpoint() string {
	return os.Getenv("AWS_S3_ENDPOINT")
}
