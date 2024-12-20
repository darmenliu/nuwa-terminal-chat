package nuwa

import (
	"context"
	"fmt"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/agents"
	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"
	"github.com/pterm/pterm"
	lcagents "github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	lcllms "github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

type NuwaAgent struct {
	ctx          context.Context
	model        lcllms.Model
	chatHistory  []lcllms.MessageContent
	systemPrompt string
	currentDir   string
	prefix       string
	catchdir     string
	scriptsdir   string
}

func NewNuwaAgent(ctx context.Context, systemPrompt string) (*NuwaAgent, error) {
	model, err := llms.GetLLMBackend(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM backend: %w", err)
	}

	return &NuwaAgent{
		ctx:          ctx,
		model:        model,
		chatHistory:  []lcllms.MessageContent{},
		systemPrompt: systemPrompt,
		currentDir:   "",
		prefix:       AgentModePrefix,
		catchdir:     NuwaCatchDir,
		scriptsdir:   NuwaScriptsDir,
	}, nil
}

func (n *NuwaAgent) Run(prompt string) error {
	return n.handleAgentMode(prompt)
}

func (n *NuwaAgent) handleAgentMode(input string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	llm, err := llms.GetLLMBackend(n.ctx)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to get LLM backend,", logger.Args("err", err.Error()))
		return err
	}

	agentTools := []tools.Tool{
		&agents.ScriptExecutor{},
	}

	agent := agents.NewTroubleshootingAgent(llm, agentTools, "output", nil)
	executor := lcagents.NewExecutor(agent)
	answer, err := chains.Run(context.Background(), executor, input)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to run agent,", logger.Args("err", err.Error()))
		return err
	}

	fmt.Println("NUWA: " + answer)
	return nil
}
