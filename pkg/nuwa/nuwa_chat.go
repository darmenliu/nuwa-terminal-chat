package nuwa

import (
	"context"
	"fmt"
	"strings"

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

	var fullResponse strings.Builder

	// 使用流式生成
	_, err := n.model.GenerateContent(ctx, n.chatHistory, lcllms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Printf("%s", chunk) // 实时输出到终端
		fullResponse.Write(chunk)
		return nil
	}))

	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// 将AI的响应添加到聊天历史
	n.chatHistory = append(n.chatHistory, lcllms.TextParts(lcllms.ChatMessageTypeAI, fullResponse.String()))

	return fullResponse.String(), nil
}

func (n *NuwaChat) Run(prompt string) error {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

	fmt.Printf("NUWA: ")
	_, err := n.Chat(n.ctx, prompt)
	if err != nil {
		logger.Error("NUWA TERMINAL: failed to generate content,", logger.Args("err", err.Error()))
		return err
	}

	fmt.Printf("\n")
	return nil
}
