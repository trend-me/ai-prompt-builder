package main

import (
	"context"
	"log/slog"

	"github.com/trend-me/ai-prompt-builder/internal/config/injector"
	"github.com/trend-me/ai-prompt-builder/internal/integration/connections"
)

func main() {
	ctx := context.Background()
	consumer, err := injector.InitializeConsumer()
	if err != nil {
		slog.Error("Error initializing consumer",
			slog.String("error", err.Error()),
		)
		return
	}

	_, _ = consumer.Consume(ctx)

	defer connections.Disconnect()
}
