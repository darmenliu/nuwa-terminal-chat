package nuwa

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/cmdexe"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/parser"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	lcllms "github.com/tmc/langchaingo/llms"
)

type NuwaTask struct {
	ctx          context.Context
	model        lcllms.Model
	chatHistory  []lcllms.MessageContent
	systemPrompt string
	currentDir   string
	prefix       string
	catchdir     string
	scriptsdir   string
}

func NewNuwaTask(ctx context.Context, systemPrompt string) (*NuwaTask, error) {
	model, err := llms.GetLLMBackend(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM backend: %w", err)
	}

	return &NuwaTask{
		ctx:          ctx,
		model:        model,
		chatHistory:  []lcllms.MessageContent{},
		systemPrompt: systemPrompt,
		currentDir:   "",
		prefix:       TaskModePrefix,
		catchdir:     NuwaCatchDir,
		scriptsdir:   NuwaScriptsDir,
	}, nil
}

func (n *NuwaTask) Run(prompt string) error {
	return n.handleTaskMode(n.ctx, prompt)
}

// handleTaskMode 处理任务模式
func (n *NuwaTask) handleTaskMode(ctx context.Context, prompt string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	rsp, err := llms.GenerateContent(ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}
	fmt.Println("NUWA: " + rsp)

	if err := n.parseScriptAndExecute(rsp); err != nil {
		logger.Error("NUWA TERMINAL: failed to parse script and execute,", logger.Args("err", err.Error()))
		return err
	}

	return nil
}

func (n *NuwaTask) parseScriptAndExecute(rsp string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	filename, content, err := n.ParseScript(rsp)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to parse script,", logger.Args("err", err.Error()))
		return err
	}

	if filename == "" {
		logger.Info("NUWA TERMINAL: empty script")
		return nil
	}

	scriptfile, err := n.prepareScriptFile(filename, content)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to prepare script file,", logger.Args("err", err.Error()))
		return err
	}

	output, err := cmdexe.ExecScriptWithOutput(scriptfile)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to execute script,", logger.Args("err", err.Error()), logger.Args("output", output))
		return err
	}

	logger.Info("NUWA TERMINAL: script output", logger.Args("output", output))

	if err := os.Remove(scriptfile); err != nil {
		logger.Error("NUWA TERMINAL: failed to remove script file,", logger.Args("err", err.Error()))
		return err
	}
	logger.Info("NUWA TERMINAL: script file removed")
	return nil
}

// prepareScriptFile prepare script file
func (n *NuwaTask) prepareScriptFile(filename, content string) (string, error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	homedir := os.Getenv("HOME")
	scriptdir := filepath.Join(homedir, n.catchdir, n.scriptsdir)
	if err := os.MkdirAll(scriptdir, os.ModePerm); err != nil {
		logger.Error("NUWA TERMINAL: failed to create script directory,", logger.Args("err", err.Error()))
		return "", err
	}

	scriptfile := filepath.Join(scriptdir, filename)
	if err := os.WriteFile(scriptfile, []byte(content), os.ModePerm); err != nil {
		logger.Error("NUWA TERMINAL: failed to write script file,", logger.Args("err", err.Error()))
		return "", err
	}

	logger.Info("NUWA TERMINAL: script file saved to " + scriptfile)
	return scriptfile, nil
}

// ParseScript parses the code from the LLM response and returns the filename and content of the first source file.
// If there are no source files found or if there is an error parsing the code, an error is returned.
func (n *NuwaTask) ParseScript(response string) (filename, content string, err error) {
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

	sources[0].FileName = uuid.New().String() + ".sh"
	sources[0].ParseFileContent()

	filename = sources[0].FileName
	content = sources[0].FileContent
	return filename, content, nil
}
