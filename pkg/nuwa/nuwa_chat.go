package nuwa

import (
	"context"
	"errors"
	"fmt"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"

	"github.com/pterm/pterm"
	lcllms "github.com/tmc/langchaingo/llms"
)

type NuwaChat struct {
	ctx          context.Context
	model        lcllms.Model
	chatHistory  []lcllms.MessageContent
	SystemPrompt string
}

func NewNuwaChat(ctx context.Context, systemPrompt string) (*NuwaChat, error) {
	model, err := llms.GetLLMBackend(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM backend: %w", err)
	}

	content := []lcllms.MessageContent{
		lcllms.TextParts(lcllms.ChatMessageTypeSystem, systemPrompt),
	}

	return &NuwaChat{
		ctx:          ctx,
		model:        model,
		chatHistory:  content,
		SystemPrompt: systemPrompt,
	}, nil
}

func (n *NuwaChat) Chat(ctx context.Context, message string) (string, error) {
	n.chatHistory = append(n.chatHistory, lcllms.TextParts(lcllms.ChatMessageTypeHuman, message))

	resp, err := n.model.GenerateContent(ctx, n.chatHistory)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return "", errors.New("empty response from model")
	}
	c1 := choices[0]
	return c1.Content, nil
}

func (n *NuwaChat) Run(prompt string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	rsp, err := n.Chat(n.ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}

	fmt.Println("NUWA: " + rsp)
	return nil
}
