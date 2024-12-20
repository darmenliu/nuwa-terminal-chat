package nuwa

import (
	"context"
	"fmt"
	"os"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/cmdexe"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/parser"
	"github.com/pterm/pterm"
	lcllms "github.com/tmc/langchaingo/llms"
)

type NuwaCmd struct {
	ctx          context.Context
	model        lcllms.Model
	chatHistory  []lcllms.MessageContent
	systemPrompt string
	currentDir   string
	prefix       string
}

func NewNuwaCmd(ctx context.Context, systemPrompt string, prefix string) (*NuwaCmd, error) {
	model, err := llms.GetLLMBackend(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM backend: %w", err)
	}

	return &NuwaCmd{
		ctx:          ctx,
		model:        model,
		chatHistory:  []lcllms.MessageContent{},
		systemPrompt: systemPrompt,
		currentDir:   "",
		prefix:       prefix,
	}, nil
}

func (n *NuwaCmd) CheckDirChanged(in string) bool {
	if in == n.currentDir {
		return false
	}
	n.SetCurrentDir(in)
	return true
}

func (n *NuwaCmd) SetCurrentDir(in string) {
	n.currentDir = in
}

func (n *NuwaCmd) Run(prompt string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	rsp, err := llms.GenerateContent(n.ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}
	fmt.Println("NUWA: " + rsp)

	cmd, err := parser.ParseCmdFromString(rsp)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to parse command,", logger.Args("err", err.Error()))
		return err
	}

	if cmd == "" {
		logger.Info("NUWA TERMINAL: empty command")
		return nil
	}

	output, err := cmdexe.ExecCommandWithOutput(cmd)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to execute command,", logger.Args("err", err.Error(), "output", output))
		return err
	}
	fmt.Println(output)

	// check current dir
	curDir, err := os.Getwd()
	if err != nil {
		logger.Warn("NUWA TERMINAL: failed to get current directory path,", logger.Args("err", err.Error()))
		return err
	}

	if n.CheckDirChanged(curDir) {
		LivePrefixState.LivePrefix = n.currentDir + n.prefix + " "
		LivePrefixState.IsEnable = true
	}

	return nil
}
