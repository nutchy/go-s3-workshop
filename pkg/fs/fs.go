package fs

import (
	"bufio"
	"context"
	"io"
	"os"
	"path"
)

type LocalStorage struct{}

func New() *LocalStorage {
	return &LocalStorage{}
}

func (f *LocalStorage) Upload(ctx context.Context, p string, body io.ReadSeeker) error {

	if err := os.MkdirAll(path.Dir(p), 0700); err != nil {
		panic(err)
	}

	// open output file
	fo, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	reader := bufio.NewReader(body)

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
	return nil
}

func (f *LocalStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	return os.Open(path)
}
