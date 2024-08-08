package main

import (
	"context"
	"github.com/trend-me/ai-prompt-builder/internal/config/connections"
	"github.com/trend-me/ai-prompt-builder/internal/config/injector"
	"log/slog"
)

func main() {
	consumer, err := injector.InitializeConsumer(context.Background())
	if err != nil {
		slog.Error("Error initializing consumer",
			slog.String("error", err.Error()),
		)
		return
	}

	consumer()

	defer connections.Disconnect()
}
