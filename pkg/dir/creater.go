package dir

import (
	"fmt"
	"os"
)

// DirectoryCreator is an interface for creating directories.
type DirectoryCreator interface {
	CreateDir(path string) error
}

// DefaultDirectoryCreator is the default implementation of DirectoryCreator.
type DefaultDirectoryCreator struct{}

func NewDefaultDirectoryCreator() DirectoryCreator {
	return &DefaultDirectoryCreator{}
}

// CreateDir creates a directory at the specified path.
func (d DefaultDirectoryCreator) CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}
