# Thumbnail service

Makes image* thumbnails according to the following logic:
>The image is scaled down to fill the given width and height while retaining the original aspect ratio and with all of the original image visible. If the requested dimensions are bigger than the original image's, the image doesn’t scale up. If the proportions of the original image do not match the given width and height, black padding is added to the image to reach the required size

_* Supported input formats: jpeg, gif, png. Output format: jpeg_

## Installation
1. Install dep (unless you already have it)
```bash
go get -u github.com/golang/dep/cmd/dep
```
2. Install service
```bash
go get github.com/Bobochka/thumbnail_service
```

## Run
Service uses **AWS S3** to store files and **Redis** for storing locks. <br>

Optional ENV params to config the service:

| PARAM | Default value | Description |
| ------ | ------ | ------ |
| AWS_S3_ENDPOINT | aws s3 url | if you want to use s3 services that is not aws | 
| REDIS_URL | redis://localhost:6379 | url of redis instance |
| PORT | 8080 | on which port server is listening |
| AWS_REGION | us-east-1 | aws region name |
| S3_BUCKET_NAME | cldnrthumbnails | S3 bucket name |

### Running with fake-s3 and local redis:
If you don't want to use real S3, you can run fake-s3 in a docker container

```bash
docker run --name my_s3 -p 4569:4569 -d lphoward/fake-s3

# add following record to /etc/hosts
127.0.0.1 <s3_bucket_name>.localhost
```
Then start service as follows:
```bash
cd $GOPATH/src/github.com/Bobochka/thumbnail_service
go build . && AWS_S3_ENDPOINT=http://localhost:4569 ./thumbnail_service
```

## Testing
Integration tests require S3 and Redis as well:

```bash
docker run --name my_s3 -p 4569:4569 -d lphoward/fake-s3

# add following record to /etc/hosts
127.0.0.1 <s3_bucket_name>.localhost
```
Run tests:
```bash
cd $GOPATH/src/github.com/Bobochka/thumbnail_service
AWS_S3_ENDPOINT=http://localhost:4569 go test ./...
```

## Endpoints 

`GET /thumbnail`

Params:

| Name | Location | Type | Description |
| ------ | ------ | ------ | ------ |
| url | query string | string | A url pointing to the origin image | 
| width | query string | int | Result thumbnail width | 
| height | query string | int | Result thumbnail width | 

Example:
```
localhost:8080/thumbnail?url=http://foo.com/sample.jpg&width=500&height=500
```
