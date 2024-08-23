package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

const promptRoadMaps = "prompt_road_maps"

type (
	UrlApiPromptRoadMapConfig func() string

	PromptRoadMapConfig struct {
		url UrlApiPromptRoadMapConfig
	}
)

func (p PromptRoadMapConfig) GetPromptRoadMap(ctx context.Context, promptRoadMapConfigName string, promptRoadMapStep int) (*models.PromptRoadMap, error) {
	slog.InfoContext(ctx, "PromptRoadMap.GetPromptRoadMap",
		slog.String("details", "process started"))

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s/%s/%d", p.url(), promptRoadMapConfigName, promptRoadMaps, promptRoadMapStep),
		nil,
	)
	if err != nil {
		return nil, err
	}

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
			return nil, exceptions.NewPromptRoadMapConfigNotFoundError(
				fmt.Sprintf("prompt_road_map with name '%s' and step '%d' not found",
					promptRoadMapConfigName, promptRoadMapStep))

		}
		return nil, exceptions.NewGetPromptRoadMapConfigError(
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
		return nil, exceptions.NewGetPromptRoadMapConfigError("error unmarshalling response", err.Error())
	}

	return &response, nil
}

func NewPromptRoadMapConfig(url UrlApiPromptRoadMapConfig) interfaces.ApiPromptRoadMapConfig {
	return &PromptRoadMapConfig{
		url: url,
	}
}
