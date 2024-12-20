package llms

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/pterm/pterm"
	lcllms "github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	lcollama "github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

type Model interface {
	// GenerateContent generates content from a prompt.
	GenerateContent(ctx context.Context, prompt string) (string, error)

	// Chat with the model.
	Chat(ctx context.Context, messages string) (string, error)

	// Set system prompt
	SetSystemPrompt(ctx context.Context, prompt string) error

	// Close the client.
	CloseBackend() error
}

// GenerateContent generates content using the specified prompt.
// It takes a context.Context object and a prompt string as input.
// The function returns the generated content as a string and an error if any.
// The function first reads the values of environment variables LLM_BACKEND, LLM_MODEL_NAME, and LLM_TEMPERATURE.
// If LLM_BACKEND is empty, it defaults to "gemini" and LLM_MODEL_NAME defaults to "gemini-1.5-pro".
// The function then creates a model based on the value of LLM_BACKEND using the specified context, model name, and temperature.
// The model is closed using the CloseBackend method before returning the generated content.
// If an error occurs during model creation or content generation, an error is returned.
func GenerateContent(ctx context.Context, prompt string) (string, error) {

	model, err := GetLLMBackend(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get LLM backend: %w", err)
	}

	modeTemperature, err := strconv.ParseFloat(os.Getenv("LLM_TEMPERATURE"), 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse LLM_TEMPERATURE: %w", err)
	}

	resp, err := lcllms.GenerateFromSinglePrompt(ctx, model, prompt, lcllms.WithTemperature(modeTemperature))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}
	return resp, nil
}

func GetLLMBackend(ctx context.Context) (lcllms.Model, error) {
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
	llmBackend := os.Getenv("LLM_BACKEND")
	modelName := os.Getenv("LLM_MODEL_NAME")
	apiKey := os.Getenv("LLM_API_KEY")

	if llmBackend == "" {
		llmBackend = "gemini"
		modelName = "gemini-1.5-pro"
	}

	if apiKey == "" {
		logger.Error("LLM_API_KEY is not set")
		return nil, errors.New("LLM_API_KEY is not set")
	}

	var model lcllms.Model
	var err error
	switch llmBackend {
	case "gemini":
		model, err = googleai.New(ctx, googleai.WithAPIKey(apiKey), googleai.WithDefaultModel(modelName))
	case "ollama":
		serverUrl := os.Getenv("OLLAMA_SERVER_URL")
		if serverUrl == "" {
			logger.Error("OLLAMA_SERVER_URL is not set")
			return nil, errors.New("OLLAMA_SERVER_URL is not set")
		}
		model, err = lcollama.New(lcollama.WithModel(modelName), lcollama.WithServerURL(serverUrl))
	case "groq":
		model, err = openai.New(
			openai.WithModel("llama3-8b-8192"),
			openai.WithBaseURL("https://api.groq.com/openai/v1"),
			openai.WithToken(apiKey),
		)
	case "deepseek":
		baseurl := os.Getenv("LLM_BASE_URL")
		if baseurl == "" {
			logger.Error("LLM_BASE_URL is not set")
			return nil, errors.New("LLM_BASE_URL is not set")
		}
		model, err = openai.New(
			openai.WithModel(modelName),
			openai.WithBaseURL(baseurl),
			openai.WithToken(apiKey),
		)
	case "claude":
		model, err = anthropic.New(
			anthropic.WithModel("claude-3-5-sonnet-20240620"),
			anthropic.WithToken(apiKey),
		)
	default:
		return nil, fmt.Errorf("unknown LLM backend: %s", llmBackend)
	}

	if err != nil {
		logger.Error(fmt.Sprintf("failed to create %s client, error:", llmBackend), logger.Args("err", err.Error()))
		return nil, err
	}

	return model, nil
}
