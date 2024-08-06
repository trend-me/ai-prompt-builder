package main

import (
	"context"
	"github.com/trend-me/ai-prompt-builder/internal/config/connections"
	"log/slog"
)

func main() {
	consumer, err := InitializeConsumer(context.Background())
	if err != nil {
		slog.Error("Error initializing consumer",
			slog.String("error", err.Error()),
		)
		return
	}

	consumer()

	defer connections.Disconnect()
}
