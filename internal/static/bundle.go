package static


import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

type fileData struct{
  path string
  root string
  data []byte
}

var (
  assets = map[string][]string{
    
      ".tml": []string{  // all .tml assets.
        
          "appliers.tml",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "appliers.tml": { // all .tml assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x92\x41\x6b\xdc\x30\x10\x85\xcf\xd6\xaf\x78\x2c\x39\xd8\x21\x68\xe9\xa1\x97\x40\x0e\xa1\x69\x21\xd0\x96\x42\x4b\xaf\x41\x2b\xcf\xda\x6a\x6d\xc9\x3b\x92\x1d\x16\xad\xfe\x7b\x91\xbc\x09\x49\xa0\xb7\x76\xd9\x8b\xd1\x37\x6f\xde\xcc\xbc\x18\x71\xd1\x2b\xdb\x0e\x84\xeb\x1b\xd4\x7e\xde\x79\xc8\xef\x81\xe5\x57\x35\x12\xde\xe1\x84\xe0\x3e\xbb\x47\xe2\x06\x29\x89\x45\x31\x6a\x51\x6d\xb7\x88\xf1\x99\x4a\xe9\xb6\xeb\x98\x3a\x15\xe8\xfe\x0e\x4c\x13\x93\x27\x1b\x3c\x42\x4f\x98\xad\x39\xcc\x04\xf5\x44\xc0\xb4\xd8\x3b\x86\x1a\x06\xd0\x92\xb1\x22\xc7\x34\xa8\x40\x2d\x82\x2b\x55\xaf\xd4\x11\x8e\x13\x49\xdc\x07\x98\x55\x33\x7f\x17\x7b\xbd\xf2\x3d\xb5\x98\xbd\xb1\x1d\x14\xc6\xf6\x3d\xfc\x3c\x4a\x51\xfd\xd5\xdd\x0d\x62\xc4\x2f\x67\xec\x4f\xc5\x46\xb5\x46\x63\x73\xbd\x79\x31\xf0\xe6\x19\xde\xe0\x54\x04\x4f\x38\xcc\x2e\x50\x9e\xbe\x11\x62\xbb\xbd\xfc\xb7\x3f\xf1\x76\x97\xf8\x98\xb7\x82\xdb\x69\x1a\x0c\x31\xfe\x43\xc7\xdc\x32\xcb\x1f\x41\xe3\xce\xb5\x86\xd6\xad\x1a\x1b\x88\xad\x1a\x30\xb8\xce\x68\x58\xd2\xe4\xbd\xe2\x63\xbe\x89\x2a\xb8\x9f\x48\x9b\xbd\xd1\xe7\xc3\x95\x87\x37\xe6\x77\xc7\xac\xae\xd5\x30\x94\x93\x4c\x13\xbb\x89\x4d\xbe\xfb\x48\xa1\x77\xad\x97\x62\x3f\x5b\x8d\x3a\xc6\x73\xec\x52\xc2\xe5\x2b\x8d\x66\x35\x57\xd3\xe2\xa1\x0f\xec\x7f\x9b\x20\xcb\x4e\x3e\xb8\x71\x34\xa1\x01\x31\x3b\x46\x14\x55\xce\xd1\xc3\xd5\xea\x26\xa7\x97\x95\xed\x08\xb4\xf8\x95\xf7\x99\xa9\xfc\xa3\x09\xba\x07\x2d\x99\x28\xa8\xbc\x53\x41\xc9\x3a\xa7\xa8\x29\x48\x8c\xe7\xd2\x8b\x87\x2b\x5c\x4c\xca\x70\x66\xe5\x37\x65\xd8\x23\x25\xad\x7c\x4e\x64\x79\x90\x3f\xce\xd9\x4b\xe9\x5a\x54\x55\xc5\x14\x66\xb6\x78\x31\x8d\x7c\x22\xbf\x94\x81\xe5\xa7\xd9\xea\xb5\xa0\xa6\xa5\x29\xdd\xc8\xb6\x29\x89\xaa\x4a\x22\xff\xcf\x12\xd6\x0c\x22\x09\xf1\x27\x00\x00\xff\xff\x41\xd7\x6f\xed\x90\x03\x00\x00"),
          path: "appliers.tml",
          root: "appliers.tml",
        },
      
    
  }
)

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
func FindDecompressedGzippedFile(path string) (io.Reader, int64, error){
	return FindFile(path, true)
}

// MustFindDecompressedGzippedFile panics if error occured, uses FindUnGzippedFile underneath.
func MustFindDecompressedGzippedFile(path string) (io.Reader, int64){
	reader, size, err := FindDecompressedGzippedFile(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindGzippedFile returns a io.Reader by seeking the giving file path if it exists.
// It returns an uncompressed file.
func FindGzippedFile(path string) (io.Reader, int64, error){
	return FindFile(path, false)
}

// MustFindGzippedFile panics if error occured, uses FindUnGzippedFile underneath.
func MustFindGzippedFile(path string) (io.Reader, int64){
	reader, size, err := FindGzippedFile(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindFile returns a io.Reader by seeking the giving file path if it exists.
func FindFile(path string, doGzip bool) (io.Reader, int64, error){
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
func MustFindFileReader(path string) (*bytes.Reader, int64){
	reader, size, err := FindFileReader(path)
	if err != nil {
		panic(err)
	}
	return reader, size
}

// FindFileReader returns a io.Reader by seeking the giving file path if it exists.
func FindFileReader(path string) (*bytes.Reader, int64, error){
  item, ok := assetFiles[path]
  if !ok {
    return nil,0, fmt.Errorf("File %q not found in file system", path)
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
func ReadFile(path string, doGzip bool) (string, error){
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
func ReadFileByte(path string, doGzip bool) ([]byte, error){
  reader, _, err := FindFile(path, doGzip)
  if err != nil {
    return nil, err
  }

  if closer, ok := reader.(io.Closer); ok {
    defer closer.Close()
  }

  var bu bytes.Buffer

  _, err = io.Copy(&bu, reader);
  if err != nil && err != io.EOF {
   return nil, fmt.Errorf("File %q failed to be read: %+q", path, err)
  }

  return bu.Bytes(), nil
}
