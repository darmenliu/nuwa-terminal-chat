package ollama

import (
	"context"
	"fmt"

	"nuwa-engineer/pkg/llms"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	Client      *ollama.Client
	Temperature float64
}

func NewOllama(ctx context.Context, model string, temperature float64) (*llms.Model, error) {

	llm, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		return nil, fmt.Errorf("Failed to create Ollama client: %w", err)
	}

	return &Ollama{
		Client:      llm,
		Temperature: temperature,
	}, nil
}

func (o *Ollama) GenerateContent(ctx context.Context, prompt string) (string, error) {
	completion, err := llm.Call(ctx, prompt,
		llms.WithTemperature(0.8),
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

func (o *Ollama) CloseBackend() error {
	return nil
}
