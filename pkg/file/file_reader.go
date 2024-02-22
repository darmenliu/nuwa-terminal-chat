package file

import (
	"os"
)

// FileReader is an interface for reading content from a file.
type FileReader interface {
	ReadFile(path string) (string, error)
}

// DefaultFileReader is the default implementation of FileReader.
type DefaultFileReader struct{}

func newDefaultFileReader() FileReader {
	return &DefaultFileReader{}
}

// ReadFile reads the content of a file and returns it as a string.
func (fr *DefaultFileReader) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Usage:
// fr := &DefaultFileReader{}
// content, err := fr.ReadFile("/path/to/file.txt")
// if err != nil {
//     // handle error
// }
// fmt.Println(content)
