package interfaces

import (
	"context"
)

type QueueAiPromptBuilderConsumer interface {
	Consume(ctx context.Context) (chan error, error)
}
