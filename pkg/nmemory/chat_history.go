package nmemory

import (
	"context"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/memory"
)

type ChatHistory struct {
	history *memory.ChatMessageHistory
}

func NewChatHistory() *ChatHistory {
	return &ChatHistory{history: memory.NewChatMessageHistory()}
}

func (h *ChatHistory) Messages(_ context.Context) ([]llms.ChatMessage, error) {
	return h.history.Messages(_)
}

func (h *ChatHistory) AddAIMessage(_ context.Context, text string) error {
	h.history.AddAIMessage(text)
	return nil
}

func (h *ChatHistory) AddUserMessage(_ context.Context, text string) error {
	h.history.AddUserMessage(text)
	return nil
}

func (h *ChatHistory) Clear(_ context.Context) error {
	h.history.Clear()
	return nil
}

func (h *ChatHistory) AddMessage(_ context.Context, message llms.ChatMessage) error {
	h.history.AddMessage(message)
	return nil
}
