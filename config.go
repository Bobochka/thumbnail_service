package main

import (
	"fmt"
	"os"

	"github.com/Bobochka/thumbnail_service/lib/downloader"
	"github.com/Bobochka/thumbnail_service/lib/locker"
	"github.com/Bobochka/thumbnail_service/lib/service"
	"github.com/Bobochka/thumbnail_service/lib/store"
)

const (
	defaultBucket   = "cldnrthumbnails"
	defaultRedisURL = "localhost:6379"
	defaultBindPort = "8080"
	awsRegion       = "us-east-1"
)

func ReadConfig() (*service.Config, error) {
	store, err := store.New(os.Getenv("AWS_S3_ENDPOINT"), awsRegion, bucketName())
	if err != nil {
		return nil, fmt.Errorf("unable to init s3 store: ", err)
	}

	locker, err := locker.New(redisURL())
	if err != nil {
		return nil, fmt.Errorf("unable to connect to redis: ", err)
	}

	return &service.Config{
		Store:      store,
		Downloader: downloader.Http{},
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
