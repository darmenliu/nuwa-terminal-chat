package workspace

import (
	"os"
)

// GetWorkspacePath retrieves the workspace path from the NUWA_WORKSPACE environment variable.
// If the environment variable is not set, an empty string is returned.
func GetWorkspacePath() string {
	return os.Getenv("NUWA_WORKSPACE")
}

// Check the dir exist
func IsWorkspaceExist(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

