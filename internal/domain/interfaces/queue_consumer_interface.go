package interfaces

import (
	"context"
)

type QueueAiPromptBuilder interface {
	Consume(ctx context.Context) (chan error, error)
}
