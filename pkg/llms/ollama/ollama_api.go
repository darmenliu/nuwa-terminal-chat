package ollama

import (
	"context"
	"fmt"

	nuwaLLM "github.com/darmenliu/nuwa-terminal-chat/pkg/llms"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Ollama struct {
	Client      *ollama.LLM
	Temperature float64
	SystemPrompt string
}

func NewOllama(ctx context.Context, model string, temperature float64) (nuwaLLM.Model, error) {

	llm, err := ollama.New(ollama.WithModel(model))
	if err != nil {
		return nil, fmt.Errorf("Failed to create Ollama client: %w", err)
	}

	return &Ollama{
		Client:      llm,
		Temperature: temperature,
		SystemPrompt: "",
	}, nil
}

func (o *Ollama) GenerateContent(ctx context.Context, prompt string) (string, error) {
	completion, err := o.Client.Call(ctx, prompt,
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

func (o *Ollama) Chat(ctx context.Context, message string) (string, error) {
	// TODO chat with AI model
	return "", nil
}

func (o *Ollama) SetSystemPrompt(ctx context.Context, prompt string) error {
	o.SystemPrompt = prompt
	return nil
}

func (o *Ollama) CloseBackend() error {
	return nil
}
