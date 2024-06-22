package gemini

import (
	"context"
	"fmt"
	"os"

	"github.com/darmenliu/nuwa-terminal-chat/pkg/llms"

	"github.com/google/generative-ai-go/genai"
	lcllms "github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"google.golang.org/api/option"
)

// Gemini is a wrapper around the Gemini API.
type Gemini struct {
	Client      *genai.Client
	Model       *genai.GenerativeModel
	google      *googleai.GoogleAI
	chatHistory []lcllms.MessageContent
	SystemPrompt string
}

// NewGemini returns a new Gemini client.
func NewGemini(ctx context.Context, modelName string, systemPrompt string) (llms.Model, error) {

	genaiKey := os.Getenv("GEMINI_API_KEY")
	if genaiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	llm, err := googleai.New(ctx, googleai.WithAPIKey(genaiKey), googleai.WithDefaultModel(modelName))
	if err != nil {
		return nil, fmt.Errorf("Failed to create GoogleAI client: %w", err)
	}

	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("Failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(modelName)

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	content := []lcllms.MessageContent{
		lcllms.TextParts(lcllms.ChatMessageTypeSystem, systemPrompt),
	}

	return &Gemini{
		Client: client,
		Model:  model,
		google: llm,
		chatHistory: content,
		SystemPrompt: systemPrompt,
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

// Chat with the model.
func (g *Gemini) Chat(ctx context.Context, message string) (string, error) {

	// Add the message to the chat history
	//prompt := g.SystemPrompt + message
	g.chatHistory = append(g.chatHistory, lcllms.TextParts(lcllms.ChatMessageTypeHuman, message))

	resp, err := g.google.GenerateContent(ctx, g.chatHistory)
	if err != nil {
		return "", fmt.Errorf("Failed to generate content: %w", err)
	}

	// Add the assistant's response to the chat history
	respchoice := resp.Choices[0]
	assistantResponse := lcllms.TextParts(lcllms.ChatMessageTypeAI, respchoice.Content)
	g.chatHistory = append(g.chatHistory, assistantResponse)
	return respchoice.Content, nil
}

// Set system prompt
func (g *Gemini) SetSystemPrompt(ctx context.Context, prompt string) error {
	g.SystemPrompt = prompt
	return nil
}

// Close the client.
func (g *Gemini) CloseBackend() error {
	return g.Client.Close()
}
