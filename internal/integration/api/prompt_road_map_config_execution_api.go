package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type (
	UrlApiPromptRoadMapConfigExecution func() string

	PromptRoadMapConfigExecution struct {
		url UrlApiPromptRoadMapConfigExecution
	}
)

func (p PromptRoadMapConfigExecution) UpdatePromptRoadMapConfigExecution(ctx context.Context, promptRoadMapConfigExecution *models.PromptRoadMapConfigExecution) error {
	slog.InfoContext(ctx, "PromptRoadMap.UpdatePromptRoadMapConfigExecution",
		slog.String("details", "process started"))

	if promptRoadMapConfigExecution.Id == nil {
		return exceptions.NewValidationError("'id' is required to update prompt_road_map_config_execution")
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/%s", p.url(), *promptRoadMapConfigExecution.Id),
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	slog.DebugContext(ctx, "PromptRoadMap.UpdatePromptRoadMapConfigExecution",
		slog.String("details", "requested"),
		slog.String("response", string(body)),
		slog.String("status", resp.Status),
	)

	if resp.StatusCode != http.StatusOK {
		return exceptions.NewUpdatePromptRoadMapConfigExecutionError(
			fmt.Sprintf("response statusCode: '%s'",
				resp.Status))
	}
	return nil
}

func NewPromptRoadMapConfigExecution(url UrlApiPromptRoadMapConfigExecution) interfaces.ApiPromptRoadMapConfigExecution {
	return &PromptRoadMapConfigExecution{
		url: url,
	}
}
