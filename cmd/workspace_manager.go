package main

import (
	"fmt"
	"log/slog"
	"nuwa-engineer/pkg/dir"
	nfile "nuwa-engineer/pkg/file"
	"nuwa-engineer/pkg/workspace"
	"os"
	"path/filepath"
)

type WorkSpaceManager interface {
	// Check workspace is created
	IsWorkspaceEixst() bool
	// CreateWorkspace creates a new workspace.
	CreateWorkspace() error
	// CreateGolangProject creates a new Golang project in the workspace.
	CreateGolangProject(projectName string) error
	// InitGolangProject initializes the Golang project.
	InitGolangProject(projectPath string, description string) error
}

type DefaultWorkSpaceManager struct {
	WorkSpacePath string
}

func NewDefaultWorkSpaceManager() WorkSpaceManager {
	return &DefaultWorkSpaceManager{WorkSpacePath: "./workspace"}
}

// Check workspace is created
func (d *DefaultWorkSpaceManager) IsWorkspaceEixst() bool {
	return workspace.IsWorkspaceExist(d.WorkSpacePath)
}

// Initialize the workspace for nuwa-engineer
func (d *DefaultWorkSpaceManager) CreateWorkspace() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	workspacePath := workspace.GetWorkspacePath()
	if workspacePath == "" {
		logger.Error("failed to get workspace path")
		return fmt.Errorf("failed to get workspace path")
	}

	d.WorkSpacePath = workspacePath

	dirCreator := dir.NewDefaultDirectoryCreator()
	err := dirCreator.CreateDir(workspacePath)
	if err != nil {
		logger.Error("failed to create workspace", err.Error())
		return fmt.Errorf("failed to create workspace: %w", err)
	}

	logger.Info("workspace created", "path", workspacePath)
	return nil
}

// Create Golang project dir in the workspace, and initialize the project
func (d *DefaultWorkSpaceManager) CreateGolangProject(projectName string) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	projectPath := filepath.Join(d.WorkSpacePath, projectName)
	logger.Info("creating project", "path", projectPath)
	if projectPath == "" {
		logger.Error("failed to get project path")
		return fmt.Errorf("failed to get project path")
	}

	dirCreator := dir.NewDefaultDirectoryCreator()
	err := dirCreator.CreateDir(projectPath)
	if err != nil {
		logger.Error("failed to create project", err.Error())
		return fmt.Errorf("failed to create project: %w", err)
	}
	logger.Info("project created", "path", projectPath)
	return nil
}

// Initialize the Golang project, crate cmd, pkg, internal, and test dirs, and create readme.md
func (d DefaultWorkSpaceManager) InitGolangProject(projectPath string, description string) error {
	// Create cmd dir
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	fullPath := filepath.Join(d.WorkSpacePath, projectPath)
	err := dir.NewDefaultDirectoryCreator().CreateDir(filepath.Join(fullPath, "cmd"))
	if err != nil {
		logger.Error("failed to create cmd dir", "err", err.Error())
		return fmt.Errorf("failed to create cmd dir: %w", err)
	}

	// Create pkg dir
	err = dir.NewDefaultDirectoryCreator().CreateDir(filepath.Join(fullPath, "pkg"))
	if err != nil {
		logger.Error("failed to create pkg dir", "err", err.Error())
		return fmt.Errorf("failed to create pkg dir: %w", err)
	}

	// Create internal dir
	err = dir.NewDefaultDirectoryCreator().CreateDir(filepath.Join(fullPath, "internal"))
	if err != nil {
		logger.Error("failed to create internal dir", "err", err.Error())
		return fmt.Errorf("failed to create internal dir: %w", err)
	}

	// Create test dir
	err = dir.NewDefaultDirectoryCreator().CreateDir(filepath.Join(fullPath, "test"))
	if err != nil {
		logger.Error("failed to create test dir", "err", err.Error())
		return fmt.Errorf("failed to create test dir: %w", err)
	}

	// Create readme.md
	readmePath := filepath.Join(fullPath, "README.md")
	err = nfile.NewFileWriter().WriteToFile(readmePath, description)
	if err != nil {
		logger.Error("failed to create README.md", "err", err.Error())
		return fmt.Errorf("failed to create readme.md: %w", err)
	}

	logger.Info("project initialized", "path", fullPath)
	return nil
}
