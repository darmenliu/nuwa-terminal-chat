package gemini

import (
	"context"
	"fmt"
	"os"

	"nuwa-engineer/pkg/llms"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Gemini is a wrapper around the Gemini API.
type Gemini struct {
	Client *genai.Client
	Model  *genai.GenerativeModel
}

// NewGemini returns a new Gemini client.
func NewGemini(ctx context.Context, modelName string) (llms.Model, error) {
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("Failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(modelName)

	return &Gemini{
		Client: client,
		Model:  model,
	}, nil
}

// Cotent to string
func (g *Gemini) ContentToString(content *genai.Content) string {
	var str string
	for _, part := range content.Parts {
		// Get interface part type, and check if it is Text
		if _, ok := part.(genai.Text); ok {
			str += string(part.(genai.Text))
		}
	}
	return str
}

// GenerateContent generates content from a prompt.
func (g *Gemini) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := g.Model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("Failed to generate content: %w", err)
	}

	// convert resp to string
	return g.ContentToString(resp.Candidates[0].Content), nil
}

// Close the client.
func (g *Gemini) CloseBackend() error {
	return g.Client.Close()
}
