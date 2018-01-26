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
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x52\xc1\x6a\xdc\x30\x10\x3d\x5b\x5f\xf1\x58\xf6\x60\x87\xa0\xa5\x87\x5e\x16\x72\x08\xa4\x85\x40\x5b\x0a\x2d\xbd\x06\xad\x3c\x6b\xab\xf5\x4a\xde\x91\xbc\x61\xd1\xea\xdf\x8b\x64\x27\x24\x81\xde\x1a\xe3\x8b\x98\xf7\xde\xbc\x99\x79\x31\x62\xdd\x2b\xdb\x0e\x84\xed\x0d\x6a\x3f\xed\x3c\xe4\x8f\xc0\xf2\x9b\x3a\x10\x3e\xe0\x82\xe0\xbe\xb8\x47\xe2\x06\x29\x89\x93\x62\xd4\xa2\xda\x6c\x10\xe3\x33\x2a\xa5\xdb\xae\x63\xea\x54\xa0\xfb\x3b\x30\x8d\x4c\x9e\x6c\xf0\x08\x3d\x61\xb2\xe6\x38\x11\xd4\x13\x02\xa6\xc5\xde\x31\xd4\x30\x80\x4e\x19\x56\xe4\x98\x06\x15\xa8\x45\x70\x85\xf5\x4a\x1d\xe1\x3c\x92\xc4\x7d\x80\x99\x35\xf3\xbb\xd8\xeb\x95\xef\xa9\xc5\xe4\x8d\xed\xa0\x70\x68\x3f\xc2\x4f\x07\x29\xaa\x7f\xba\xbb\x41\x8c\xf8\xed\x8c\xfd\xa5\xd8\xa8\xd6\x68\xac\xb6\xab\x17\x03\xaf\x9e\xc1\x2b\x5c\x8a\xe0\x05\xc7\xc9\x05\xca\xd3\x37\x42\x6c\x36\x57\xff\xf7\x13\x6f\x77\x89\x4f\x79\x2b\xb8\x1d\xc7\xc1\x10\xe3\x1d\x3a\xe6\x96\x59\xfe\x0c\x3a\xec\x5c\x6b\x68\xde\xaa\xb1\x81\xd8\xaa\x01\x83\xeb\x8c\x86\x25\x4d\xde\x2b\x3e\xe7\x9b\xa8\x02\xf7\x23\x69\xb3\x37\x7a\x39\x5c\x29\xbc\x31\xbf\x3b\x67\x75\xad\x86\xa1\x9c\x64\x1c\xd9\x8d\x6c\xf2\xdd\x0f\x14\x7a\xd7\x7a\x29\xf6\x93\xd5\xa8\x63\x5c\x62\x97\x12\xae\x5e\x69\x34\xb3\xb9\x9a\x4e\x1e\x52\x4a\x7d\x64\xff\xc7\x04\x59\xd6\xd2\x80\x98\x1d\x23\x8a\x2a\x87\xe8\xe1\x7a\xb6\x92\xa3\xcb\xca\x76\x84\x4c\x8a\xa2\xaa\xfc\xa3\x09\xba\x07\x9d\x72\xa9\x60\x66\x81\x3b\x15\x94\xac\x73\x80\x9a\x82\x8b\x71\x21\xae\x1f\xae\xb1\x1e\x95\xe1\x4c\x90\xdf\x95\x61\x8f\x94\xb4\xf2\x39\x8c\xa5\x20\x7f\x2e\xb1\x4b\x69\x2b\xaa\xaa\x62\x0a\x13\x5b\xbc\x18\x44\x3e\x21\xbf\x96\x59\xe5\xe7\xc9\xea\x99\x50\xd3\xa9\x29\xdd\xc8\xb6\x29\x89\xaa\x4a\x22\xff\x8b\x84\x35\x83\x48\x42\x08\xf1\x37\x00\x00\xff\xff\x8c\xf1\x4d\x81\x8d\x03\x00\x00"),
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
