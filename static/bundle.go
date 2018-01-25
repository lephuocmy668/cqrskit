package static

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

var (
	assets     = map[string][]string{}
	assetFiles = map[string]fileData{}
)

type fileData struct {
	path string
	root string
	data []byte
}

//==============================================================================

// FilesFor returns all files that use the provided extension, returning a
// empty/nil slice if none is found.
func FilesFor(ext string) []string {
	return assets[ext]
}

// MustFindFile calls FindFile to retrieve file reader with path else panics.
func MustFindFile(path string, doGzip bool) (io.Reader, int64) {
	reader, size, err := FindFile(path, doGzip)
	if err != nil {
		panic(err)
	}

	return reader, size
}

// FindDecompressedGzippedFile returns a io.Reader by seeking the giving file path if it exists.
// It returns an uncompressed file.
func FindDecompressedGzippedFile(path string) (io.Reader, int64, error) {
	return FindFile(path, true)
}

// MustFindDecompressedGzippedFile panics if error occured, uses FindUnGzippedFile underneath.
func MustFindDecompressedGzippedFile(path string) (io.Reader, int64) {
	reader, size, err := FindDecompressedGzippedFile(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindGzippedFile returns a io.Reader by seeking the giving file path if it exists.
// It returns an uncompressed file.
func FindGzippedFile(path string) (io.Reader, int64, error) {
	return FindFile(path, false)
}

// MustFindGzippedFile panics if error occured, uses FindUnGzippedFile underneath.
func MustFindGzippedFile(path string) (io.Reader, int64) {
	reader, size, err := FindGzippedFile(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindFile returns a io.Reader by seeking the giving file path if it exists.
func FindFile(path string, doGzip bool) (io.Reader, int64, error) {
	reader, size, err := FindFileReader(path)
	if err != nil {
		return nil, size, err
	}

	if !doGzip {
		return reader, size, nil
	}

	gzr, err := gzip.NewReader(reader)
	return gzr, size, err
}

// MustFindFileReader returns bytes.Reader for path else panics.
func MustFindFileReader(path string) (*bytes.Reader, int64) {
	reader, size, err := FindFileReader(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindFileReader returns a io.Reader by seeking the giving file path if it exists.
func FindFileReader(path string) (*bytes.Reader, int64, error) {
	item, ok := assetFiles[path]
	if !ok {
		return nil, 0, fmt.Errorf("File %q not found in file system", path)
	}

	return bytes.NewReader(item.data), int64(len(item.data)), nil
}

// MustReadFile calls ReadFile to retrieve file content with path else panics.
func MustReadFile(path string, doGzip bool) string {
	body, err := ReadFile(path, doGzip)
	if err != nil {
		panic(err)
	}

	return body
}

// ReadFile attempts to return the underline data associated with the given path
// if it exists else returns an error.
func ReadFile(path string, doGzip bool) (string, error) {
	body, err := ReadFileByte(path, doGzip)
	return string(body), err
}

// MustReadFileByte calls ReadFile to retrieve file content with path else panics.
func MustReadFileByte(path string, doGzip bool) []byte {
	body, err := ReadFileByte(path, doGzip)
	if err != nil {
		panic(err)
	}

	return body
}

// ReadFileByte attempts to return the underline data associated with the given path
// if it exists else returns an error.
func ReadFileByte(path string, doGzip bool) ([]byte, error) {
	reader, _, err := FindFile(path, doGzip)
	if err != nil {
		return nil, err
	}

	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	var bu bytes.Buffer

	_, err = io.Copy(&bu, reader)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("File %q failed to be read: %+q", path, err)
	}

	return bu.Bytes(), nil
}
