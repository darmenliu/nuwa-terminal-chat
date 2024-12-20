package nuwa

import (
	"context"
	"fmt"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"
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

	if err := parseScriptAndExecute(rsp); err != nil {
		logger.Error("NUWA TERMINAL: failed to parse script and execute,", logger.Args("err", err.Error()))
		return err
	}

	return nil
}
