package file

import (
	"os"
)

type FileWriter interface {
	WriteToFile(path string, content string) error
}

type FileWriterImpl struct{}

func NewFileWriter() FileWriter {
	return &FileWriterImpl{}
}

func (f *FileWriterImpl) WriteToFile(path string, content string) error {
	// Open the file with write-only mode and create it if it doesn't exist
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
