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
        
          "command-event-pairs.tml",
        
          "read-write-repository.tml",
        
      },
    
  }

  assetFiles = map[string]fileData{
    
      
        "command-event-pairs.tml": { // all .tml assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x91\xc1\x6a\xf3\x30\x10\x84\xcf\xbf\x9e\x62\xc8\xc9\x86\xa0\xf0\x3f\x42\x29\x85\x9e\xd2\x43\x9f\x40\xb1\xd6\xb1\x8a\xb4\x32\xd2\xc6\x25\x08\xbd\x7b\x51\x52\xda\x4b\x9c\xf6\xd0\x8b\xc1\x68\x77\x76\xe6\x1b\xb5\xdb\xe1\x61\x9e\xfd\xb9\x94\xa7\x85\x58\xf4\xcb\xe1\x8d\x06\xd1\x7b\x13\xe8\xf2\xa9\x15\x66\x9e\xbd\xa3\x8c\xa3\x5b\x1c\x1f\x41\x6d\x0e\x12\x21\x93\xcb\x70\x9c\xc5\xf0\x40\x88\x23\x4a\xd1\xaf\x92\x6e\x48\x68\x35\x9e\x78\x40\x57\x4a\x3e\x1d\x32\x6e\x4e\xe1\x7f\xad\xeb\x0a\xfd\x4f\x2e\x3b\x5a\xda\xf6\xda\x73\x0f\x4a\x29\x26\x14\xf5\x2f\x91\x9c\x12\x5f\xff\xb3\xde\xd3\x7b\xb7\xe1\x28\x70\x61\xf6\x14\x88\x85\x2c\xce\x24\x9b\x5e\x55\xd5\xe8\x3c\x1b\xb6\x9e\x4a\xd1\x8f\x31\x04\xc3\xf6\x16\xa1\xaf\xdd\x0c\x99\x08\x81\x64\x8a\x16\x79\x30\xe3\x18\xbd\x6d\xd0\x64\x32\x82\x23\x31\x25\x23\x94\x61\xb0\x9e\xa4\x1d\xbd\x32\x1e\x63\xba\xe8\xdd\x3d\xfe\x17\x6c\x7f\x91\xb1\x1b\x82\xbd\x6f\xa4\x6f\x1e\x56\x0b\xd8\x5e\x81\xf7\xad\x81\xc5\xa4\xcf\x84\x77\x16\xbe\x8b\x6a\x03\x5b\xb0\xf3\xaa\xaa\x8f\x00\x00\x00\xff\xff\xe0\x96\x6a\xcc\xb0\x02\x00\x00"),
          path: "command-event-pairs.tml",
          root: "command-event-pairs.tml",
        },
      
        "read-write-repository.tml": { // all .tml assets.
          data: []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x92\x31\x8f\x13\x31\x10\x85\xeb\xf8\x57\x3c\xa5\x4a\xae\xd8\x15\x05\x0d\x12\xc5\x49\x20\x71\x0d\x20\x21\xd1\x20\x8a\x89\x77\xb2\xf1\xe1\xb5\x9d\xb1\x1d\x29\xda\xf3\x7f\x47\xf6\x85\x20\x24\x96\x50\xc0\x16\x2b\x8d\xec\x79\xef\x1b\xbf\x39\x91\x60\xa3\x56\x7d\x8f\x79\xee\x3e\x25\xe9\x3e\xec\x1e\x59\xa7\xee\x3d\x4d\xdc\x7e\xa5\xdc\x8f\xa3\xf0\x48\x89\x1f\xde\x40\x38\x08\x47\x76\x29\x22\x1d\x18\xd9\x99\x63\x66\xd0\x8f\x1b\x30\x03\xf6\x5e\x40\xd6\x82\x4f\xf5\x5a\x93\x16\xb6\x94\x78\x40\xf2\xad\x6b\xd1\x09\xe9\x1c\xb8\xc3\x43\x82\x79\xd6\xaf\x75\x3d\xc2\x81\xe2\x81\x07\xe4\x68\xdc\x08\xc2\x34\xbc\x44\xcc\x53\xa7\x56\x7f\x45\xfd\x1a\xf3\x8c\x47\x6f\xdc\x67\x12\x43\x83\xd1\x58\xbf\x5a\xe3\xb7\x8d\x58\x5f\x1b\xd7\x78\x6a\x46\x4f\x38\x66\x9f\x18\xa5\xa8\xad\x52\x7d\x7f\xf7\x6f\x3f\xd5\xf7\x78\x5b\xdf\x0a\xef\xc8\x0d\x96\xe5\x3f\x58\xa8\x3f\xe5\xdb\xcc\x23\xcc\x14\x2c\x4f\xd7\x68\x1d\x6b\x8e\x91\xe4\x0c\xeb\x47\xa3\x6b\x76\x14\x82\x3d\x83\xaa\x58\x64\x31\x1c\xe1\xf7\x97\x98\xdb\x31\x46\x73\xaa\xf9\xdc\xc8\x57\xd5\xff\x4d\x9c\x98\x24\xeb\x34\xab\xd5\xa5\xfe\xf2\x55\x1f\x25\x76\xad\x52\xa5\x4d\x74\xdf\x78\x78\xda\xf9\xa1\xc2\x54\x6a\xe3\x12\x8b\x23\x7b\x81\xfe\x39\xc4\x15\x3f\x06\xd6\x66\x6f\xf4\x2f\xe0\xcb\xc4\xbb\x73\x75\xd2\x64\x6d\xdb\xbc\x10\xc4\x07\x31\x75\xd5\x27\x4e\x07\x3f\xc4\x4e\xed\xb3\xd3\xd8\x70\x3c\xdd\x9a\x69\xfb\x4c\xbc\x99\xe7\x98\x77\x71\x61\xff\x5e\x94\x82\xbb\x8b\xd0\x47\xd2\xdf\x68\xe4\x52\xba\x45\xe5\x2d\x58\xc4\x0b\x66\xa5\x56\xc2\x29\x8b\x83\x33\xb6\x3e\x90\x52\xdf\x03\x00\x00\xff\xff\x91\x49\xd9\x30\xdc\x03\x00\x00"),
          path: "read-write-repository.tml",
          root: "read-write-repository.tml",
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
