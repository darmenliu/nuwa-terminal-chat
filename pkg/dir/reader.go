package dir

import (
	"os"
)

// FileReader is an interface for reading files and directories.
type FileReader interface {
	ReadDir(folderPath string) ([]string, error)
}

// DirReader implements the FileReader interface.
type DirReader struct{}

// NewDirReader creates a new DirReader.
func NewDirReader() FileReader {
	return &DirReader{}
}

// ReadDir reads the files and directories in the specified folder.
// It returns a slice of file names and an error if any.
func (dr DirReader) ReadDir(folderPath string) ([]string, error) {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}
