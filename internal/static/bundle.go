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
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x93\x41\x6b\xdc\x30\x10\x85\xcf\xd6\xaf\x78\x2c\x39\xec\x86\xe0\xa5\x87\x5e\x02\x39\x04\xd2\x42\xa0\x2d\x85\x96\x5e\x4a\x09\x8a\x3c\xeb\x55\x2b\x4b\xca\x48\x76\x58\x14\xfd\xf7\x22\xad\xb3\xcd\xb6\xdb\x5b\x6b\x8c\x41\xd6\x9b\x79\xdf\xcc\x48\x29\xe1\x6c\x2b\x6d\x67\x08\x97\x57\x58\x86\xf1\x3e\xa0\xfd\x14\xb9\xfd\x20\x07\xc2\x2b\x3c\x21\xba\x77\xee\x91\x78\x85\x9c\xc5\x24\x19\x4b\xd1\xac\xd7\x48\xe9\xa0\xca\xf9\xba\xef\x99\x7a\x19\xe9\xf6\x06\x4c\x9e\x29\x90\x8d\x01\x71\x4b\x18\xad\x7e\x18\x09\xf2\x59\x01\xdd\x61\xe3\x18\xd2\x18\xd0\x54\x64\x35\x1d\x93\x91\x91\x3a\x44\x57\xa3\x8e\xb2\x23\xee\x3c\xb5\xb8\x8d\xd0\xfb\x9c\x65\x5d\xf1\xb6\x32\x6c\xa9\xc3\x18\xb4\xed\x21\x31\x74\xaf\x11\xc6\xa1\x15\xcd\x5f\xe9\xae\x90\x12\xbe\x3b\x6d\xbf\x48\xd6\xb2\xd3\x0a\x8b\xcb\xc5\x8b\x82\x17\x07\xf1\x02\x4f\x35\xe1\x13\x1e\x46\x17\xa9\x54\xbf\x12\x62\xbd\x3e\xff\xb7\x8f\xf8\xbd\x97\x78\x53\xba\x82\x6b\xef\x8d\x26\xc6\x7f\x70\xfc\xc3\xb2\x3a\x06\xe8\xc1\x1b\x1a\x0e\x93\xb3\xa4\x28\x04\xc9\x3b\x18\xd7\x6b\x55\x46\x23\xbd\x37\x3b\xc8\x92\x20\x10\x6b\x0a\x70\x9b\x79\x8a\x75\x1b\xbd\x9e\xca\x28\x4e\x8c\x4f\x94\xef\x49\xdb\x10\x79\x54\x31\x89\x66\x5e\x7f\xfd\xa6\x1e\x38\xfc\xd0\xb1\xad\x3f\x44\xae\xc0\xd7\xd5\x9a\x86\x7b\xd7\x15\xdf\x02\xa8\x6d\x24\xb6\xd2\xcc\x7c\xbf\x78\x0f\xa4\xc1\x93\xd2\x1b\xad\x8e\x18\x8f\xe1\xee\x77\x25\xbb\x92\xc6\xd4\x33\xe4\x3d\x3b\xcf\xba\x1c\xd4\x81\xe2\xd6\x75\xa1\x15\x9b\xd1\x2a\x2c\x29\x4c\xa7\xf0\x57\x7b\xb2\x65\x4a\xf3\x35\xca\x19\xe7\x47\xba\x15\x88\xd9\x31\x92\x68\xca\xc1\xbf\xbb\xd8\xd3\x94\xeb\xc6\xd2\xf6\x04\x0a\x53\x3b\xd7\x9e\x44\xd3\x84\x47\x1d\xd5\x16\x34\x15\x45\x95\xee\x77\x6f\x64\x94\xed\xb2\xb4\x71\x55\x75\x29\xcd\xf1\x67\x77\x17\x38\xf3\x52\x73\x09\x68\x3f\x4a\xcd\x01\x39\x2b\x19\x4a\xbf\xeb\x46\xfb\x79\xbe\x31\x39\x5f\x8a\xa6\x69\x98\xe2\xc8\x16\x2f\x98\xdb\x67\xe5\xfb\x5a\x75\xfb\x76\xb4\x6a\x1f\xb0\xa4\x69\x55\xdd\xc8\x76\x39\x8b\xa6\xc9\xa2\xbc\x73\x0a\xab\x4d\x19\x90\x10\x3f\x03\x00\x00\xff\xff\xbc\x84\x35\xab\x48\x04\x00\x00"),
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
