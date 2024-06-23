# go-s3-workshop

Demonstrate to upload file to Object storage using AWS SDK

## How to run

1. Start docker compose
```bash
$ docker compose -f docker/docker-compose.yaml up 
```
- MinIO Console: You can access the MinIO web console at http://localhost:9001
- MinIO Server: You can access the MinIO server at http://localhost:9000.

2. Create the local bucket using MinIO
- login to MinIO Console
- On the side bar, Select `Object Browser` -> click the link `Create a Bucket` -> Fill the bucket name e.g. `app-bucket`
- Create `.env` file from `.env.example` then replace the required values

5. Build go app
```sh
$ go build
```

6. Run go application
```sh
./go-s3-workshop -fn fs_stdin -k outbound/01_output.txt < inbound/01_input.txt
```
```sh
./go-s3-workshop -fn fs_file  -f inbound/02_input.txt -k outbound/02_output.txt
```
```sh
./go-s3-workshop -fn s3_stdin -k outbound/01_output.txt -d 10m < inbound/01_input.txt
```

```sh
./go-s3-workshop -fn s3_file -f inbound/02_input.txt -k outbound/02_output.txt -d 10m
```
7. Teardown docker composes
```sh
docker compose down
# OR [Ctrl + C]
```

### External Library
```
$ go get github.com/aws/aws-sdk-go
$ go get github.com/joho/godotenv
```