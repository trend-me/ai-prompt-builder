package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
	"io"
	"log/slog"
	"net/http"
)

type (
	UrlApiPromptRoadMap func() string

	PromptRoadMap struct {
		url UrlApiPromptRoadMap
	}
)

func (p PromptRoadMap) GetPromptRoadMap(ctx context.Context, promptRoadMapId string) (*models.PromptRoadMap, error) {
	slog.InfoContext(ctx, "PromptRoadMap.GetPromptRoadMap",
		slog.String("details", "process started"))

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", p.url(), promptRoadMapId),
		nil,
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, exceptions.NewPromptRoadMapNotFoundError(
				fmt.Sprintf("prompt_road_map '%s' not found",
					promptRoadMapId))

		}
		return nil, exceptions.NewGetPromptRoadMapError(
			fmt.Sprintf("response with statusCode: '%s'",
				resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response models.PromptRoadMap
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, exceptions.NewGetPromptRoadMapError("error unmarshalling response", err.Error())
	}

	return &response, nil
}

func (p PromptRoadMap) UpdatePromptRoadMapConfigExecution(ctx context.Context, promptRoadMapConfigExecution *models.PromptRoadMapConfigExecution) error {
	slog.InfoContext(ctx, "PromptRoadMap.GetPromptRoadMap",
		slog.String("details", "process started"))

	if promptRoadMapConfigExecution.Id == nil {
		return exceptions.NewValidationError("'id' is required to update prompt_road_map_config_execution")
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", p.url(), *promptRoadMapConfigExecution.Id),
		nil,
	)

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
	slog.DebugContext(ctx, "PromptRoadMap.GetPromptRoadMap",
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

func NewPromptRoadMap(url UrlApiPromptRoadMap) interfaces.ApiPromptRoadMap {
	return &PromptRoadMap{
		url: url,
	}
}
