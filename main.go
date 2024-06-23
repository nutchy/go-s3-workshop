package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"go-s3-workshop/pkg/fs"
	myS3 "go-s3-workshop/pkg/s3"

	"github.com/joho/godotenv"
)

type Storager interface {
	Upload(ctx context.Context, path string, body io.ReadSeeker) error
	Download(ctx context.Context, path string) (io.ReadCloser, error)
}

// Uploads a file to S3 given a bucket and object key. Also takes a duration
// value to terminate the update if it doesn't complete within that time.
//
// The AWS Region needs to be provided in the AWS shared config or on the
// environment variable as `AWS_REGION`. Credentials also must be provided
// Will default to shared config file, but can load from environment if provided.
//
// Usage:
//
//	# Upload myfile.txt to myBucket/myKey. Must complete within 10 minutes or will fail
//	go run withContext.go -k myKey -d 10m < myfile.txt
func main() {
	var key, fn, fileName string
	var timeout time.Duration

	flag.StringVar(&fn, "fn", "", "Function name.")
	flag.StringVar(&key, "k", "", "Object key name.")
	flag.StringVar(&fileName, "f", "", "File name.")
	flag.DurationVar(&timeout, "d", 0, "Upload timeout.")
	flag.Parse()

	envFile, err := godotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	if cancelFn != nil {
		defer cancelFn()
	}

	s3Storage := myS3.NewWithConfig(myS3.S3Config{
		Region:   "us-east-1",
		Endpoint: envFile["S3_ENDPOINT"],
		Bucket:   envFile["S3_BUCKET_NAME"],
		ID:       envFile["S3_ACCESS_KEY"],
		Secret:   envFile["S3_SECRET_KEY"],
		Timeout:  timeout,
	})

	switch fn {
	case "fs_stdin":
		fs := fs.New()
		upload(ctx, fs, key, os.Stdin)
		download(ctx, fs, key)
	case "fs_file":
		fs := fs.New()
		// Open file for reading
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		// Read the contents of the file into a buffer
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			return
		}

		// bytes reader implementes ReadSeeker interface
		upload(ctx, fs, key, bytes.NewReader(buf.Bytes()))
		download(ctx, fs, key)
	case "s3_stdin":
		upload(ctx, s3Storage, key, os.Stdin)
		download(ctx, s3Storage, key)
	case "s3_file":
		// Open file for reading
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		// Read the contents of the file into a buffer
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, file); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			return
		}

		// bytes reader implementes ReadSeeker interface
		upload(ctx, s3Storage, key, bytes.NewReader(buf.Bytes()))
		download(ctx, s3Storage, key)
	default:
		fmt.Println("Error: missing function name")
		os.Exit(1)
	}

	fmt.Println("=== DONE ===")
}

func upload(ctx context.Context, storage Storager, path string, body io.ReadSeeker) {
	if err := storage.Upload(ctx, path, body); err != nil {
		panic(err)
	}
}

func download(ctx context.Context, storage Storager, path string) {
	reader, err := storage.Download(ctx, path)
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
