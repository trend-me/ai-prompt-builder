package api

import (
	"bytes"
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
	UrlApiValidation func() string

	Validation struct {
		url UrlApiValidation
	}
)

func (v Validation) ExecutePayloadValidator(ctx context.Context, payloadValidatorId string, payload []byte) (*models.PayloadValidatorExecutionResponse, error) {
	slog.InfoContext(ctx, "Validation.ExecutePayloadValidator",
		slog.String("details", "process started"))

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", v.url(), payloadValidatorId),
		bytes.NewReader(payload),
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
			return nil, exceptions.NewPayloadValidatorNotFoundError(
				fmt.Sprintf("payload_validator '%s' not found",
					payloadValidatorId))

		}
		return nil, exceptions.NewPayloadValidatorError(
			fmt.Sprintf("response with statusCode: '%s'",
				resp.Status))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, exceptions.NewPayloadValidatorError("error unmarshalling response ", err.Error())
	}

	slog.DebugContext(ctx, "Validation.ExecutePayloadValidator",
		slog.String("details", "requested"),
		slog.String("response", string(body)),
		slog.String("status", resp.Status),
	)

	var response models.PayloadValidatorExecutionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, exceptions.NewPayloadValidatorError("error unmarshalling response ", err.Error())
	}

	return &response, nil
}

func NewValidation(url UrlApiValidation) interfaces.ApiValidation {
	return &Validation{
		url: url,
	}
}
