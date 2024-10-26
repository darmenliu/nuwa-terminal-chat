package nmemory

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
)

type ChatHistory struct {
	history *memory.ChatMessageHistory
}

// NewChatHistory creates a new ChatHistory instance
func NewChatHistory() *ChatHistory {
	return &ChatHistory{
		history: memory.NewChatMessageHistory(),
	}
}

// AddUserMessage adds a user message to the chat history
func (ch *ChatHistory) AddUserMessage(ctx context.Context, content string) error {
	return ch.history.AddUserMessage(ctx, content)
}

// AddAIMessage adds an AI message to the chat history
func (ch *ChatHistory) AddAIMessage(ctx context.Context, content string) error {
	return ch.history.AddAIMessage(ctx, content)
}

// GetMessages returns all messages in the chat history
func (ch *ChatHistory) GetMessages(ctx context.Context) ([]llms.ChatMessage, error) {
	return ch.history.Messages(ctx)
}

// Clear clears all messages from the chat history
func (ch *ChatHistory) Clear(ctx context.Context) error {
	return ch.history.Clear(ctx)
}
