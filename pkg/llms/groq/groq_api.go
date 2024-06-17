package groq

import (
	"context"
	"fmt"

	nuwaLLM "nuwa-engineer/pkg/llms"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Groq struct {
	Client      *openai.LLM
	Temperature float64
}

func NewGroq(ctx context.Context, model string, temperature float64) (nuwaLLM.Model, error) {
	llm, err := openai.New(
		openai.WithModel(model),
		openai.WithBaseURL("https://api.groq.com/openai/v1"),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Groq client: %w", err)
	}
	return &Groq{
		Client:      llm,
		Temperature: temperature,
	}, nil
}

func (o *Groq) GenerateContent(ctx context.Context, prompt string) (string, error) {
	completion, err := llms.GenerateFromSinglePrompt(ctx,
		o.Client,
		prompt,
		llms.WithTemperature(0.8),
		llms.WithMaxTokens(4096),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	if err != nil {
		return "", err
	}

	return completion, nil
}

func (o *Groq) Chat(ctx context.Context, message string) (string, error) {
	// TODO chat with AI model
	return "", nil
}

func (o *Groq) CloseBackend() error {
	return nil
}
