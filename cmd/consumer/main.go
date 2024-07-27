package main

import (
	"github.com/trend-me/ai-prompt-builder/internal/config/connections"
	"log/slog"
)

func main() {
	consumer, err := InitializeConsumer()
	if err != nil {
		slog.Error("Error initializing consumer",
			slog.String("error", err.Error()),
		)
		return
	}

	consumer()

	defer connections.Disconnect()
}
