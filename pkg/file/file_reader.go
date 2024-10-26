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

// NewDefaultFileReader creates a new instance of DefaultFileReader.
// This function is exported so it can be used by other packages if needed.
func NewDefaultFileReader() FileReader {
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
// fr := file.NewDefaultFileReader()
// content, err := fr.ReadFile("/path/to/file.txt")
// if err != nil {
//     // handle error
// }
// fmt.Println(content)
