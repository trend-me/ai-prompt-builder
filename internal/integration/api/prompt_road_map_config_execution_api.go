package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
)

type (
	UrlApiPromptRoadMapConfigExecution func() string

	PromptRoadMapConfigExecution struct {
		url UrlApiPromptRoadMapConfigExecution
	}
)

func (p PromptRoadMapConfigExecution) UpdateStepInExecutionById(ctx context.Context, id string, stepInExecution int) error {
	slog.InfoContext(ctx, "PromptRoadMap.UpdateStepInExecutionById",
		slog.String("details", "process started"),
		slog.String("id", id),
		slog.Int("stepInExecution", stepInExecution),
	)

	reqBodyStr := fmt.Sprintf(`{"step_in_execution": %d}`, stepInExecution)

	reqBody := io.NopCloser(
		strings.NewReader(reqBodyStr),
	)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/%s", p.url(), id),
		reqBody,
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
	slog.DebugContext(ctx, "PromptRoadMap.UpdateStepInExecutionById",
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
