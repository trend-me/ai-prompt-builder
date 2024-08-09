package main

import (
	"context"
	"github.com/trend-me/ai-prompt-builder/internal/config/connections"
	"github.com/trend-me/ai-prompt-builder/internal/config/injector"
	"log/slog"
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

	_, _ = consumer(ctx)

	defer connections.Disconnect()
}
