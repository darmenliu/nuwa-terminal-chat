package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/parser"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"github.com/tmc/langchaingo/tools"
)

const (
	Catchdir   = ".nuwa-terminal"
	ScriptsDir = "scripts"
)

type ScriptCodeParser struct {
}

var _ tools.Tool = &ScriptCodeParser{}

// Description returns a string describing the ScriptCodeParser tool.
func (e *ScriptCodeParser) Description() string {
	return `Useful for parsing the code snippets. 
	The input to this tool should be a code snippet`
}

// Name returns the name of the tool.
func (e *ScriptCodeParser) Name() string {
	return "ScriptCodeParser"
}

func (e *ScriptCodeParser) Call(ctx context.Context, input string) (string, error) {
	return e.ParseScriptAndSave(input)
}

func (e *ScriptCodeParser) ParseScriptAndSave(input string) (string, error) {
	filename, content, err := ParseScript(input)
	if err != nil {
		return "", err
	}

	scriptfile, err := SaveScript(filename, content)
	if err != nil {
		return "", err
	}

	return scriptfile, nil
}

func ParseScript(response string) (filename, content string, err error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	parser := parser.NewGoCodeParser()
	sources, err := parser.ParseCode(response)
	if err != nil {
		logger.Error("Failed to parse code from LLM response, error:", logger.Args("err", err.Error()))
		return "", "", err
	}

	if len(sources) == 0 {
		logger.Error("No source files found in LLM response")
		return "", "", fmt.Errorf("no source files found")
	}

	sources[0].ParseFileContent()

	filename = uuid.New().String() + ".sh"
	content = sources[0].FileContent
	return filename, content, nil
}

func SaveScript(filename, content string) (string, error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	homedir := os.Getenv("HOME")
	scriptdir := filepath.Join(homedir, Catchdir, ScriptsDir)
	err := os.MkdirAll(scriptdir, os.ModePerm)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to create script directory,", logger.Args("err", err.Error()))
		return "", err
	}

	scriptfile := filepath.Join(scriptdir, filename)
	err = os.WriteFile(scriptfile, []byte(content), os.ModePerm)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to write script file,", logger.Args("err", err.Error()))
		return "", err
	}

	return scriptfile, nil
}
