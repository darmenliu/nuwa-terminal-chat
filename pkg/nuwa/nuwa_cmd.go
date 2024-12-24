package nuwa

import (
	"context"
	"fmt"

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
}

func NewNuwaCmd(ctx context.Context, systemPrompt string) (*NuwaCmd, error) {
	model, err := llms.GetLLMBackend(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM backend: %w", err)
	}

	return &NuwaCmd{
		ctx:          ctx,
		model:        model,
		chatHistory:  []lcllms.MessageContent{},
		systemPrompt: systemPrompt,
	}, nil
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
	return nil
}
