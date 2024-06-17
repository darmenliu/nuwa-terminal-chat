package llms

import (
	"context"
)

type Model interface {
	// GenerateContent generates content from a prompt.
	GenerateContent(ctx context.Context, prompt string) (string, error)

	// Chat with the model.
	Chat(ctx context.Context, messages string) (string, error)

	// Close the client.
	CloseBackend() error
}
