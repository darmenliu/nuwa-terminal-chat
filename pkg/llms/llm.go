package llms

import (
	"context"
)

type Model interface {
	// GenerateContent generates content from a prompt.
	GenerateContent(ctx context.Context, prompt string) (string, error)
	// Close the client.
	CloseBackend() error
}
