package usecases

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/builders"
	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type UseCase struct {
	apiPromptRoadMapConfig          interfaces.ApiPromptRoadMapConfig
	apiPromptRoadMapConfigExecution interfaces.ApiPromptRoadMapConfigExecution
	apiValidation                   interfaces.ApiValidation
	queueAiRequester                interfaces.QueueAiRequester
}

func (u UseCase) Handle(ctx context.Context, request *models.Request) error {
	slog.InfoContext(ctx, "useCase.Handle",
		slog.String("details", "process started"))

	promptRoadMap, err := u.apiPromptRoadMapConfig.GetPromptRoadMap(ctx, request.PromptRoadMapConfigName, request.PromptRoadMapStep)
	if err != nil {
		return err
	}

	err = u.apiPromptRoadMapConfigExecution.UpdatePromptRoadMapConfigExecution(ctx, &models.PromptRoadMapConfigExecution{
		Id:              &request.PromptRoadMapConfigExecutionId,
		StepInExecution: promptRoadMap.Step,
	})
	if err != nil {
		return err
	}

	err = u.validateMetadata(ctx, promptRoadMap, request)
	if err != nil {
		return err
	}

	prompt, err := builders.BuildPrompt(request, promptRoadMap)
	if err != nil {
		return err
	}

	err = u.queueAiRequester.Publish(ctx, prompt, request)
	if err != nil {
		return err
	}

	slog.DebugContext(ctx, "useCase.Handle",
		slog.String("details", "process finished"))
	return nil
}

func (u UseCase) validateMetadata(ctx context.Context, promptRoadMap *models.PromptRoadMap, request *models.Request) error {
	if *promptRoadMap.Step > 1 {
		payload, err := json.Marshal(request.Metadata)
		if err != nil {
			return exceptions.NewValidationError(err.Error())
		}

		payloadValidationExecutionResponse, err := u.apiValidation.ExecutePayloadValidator(ctx, *promptRoadMap.MetadataValidationName, payload)
		if err != nil {
			return err
		}

		bPayloadValidationExecutionResponse, _ := json.Marshal(payloadValidationExecutionResponse)
		slog.InfoContext(ctx, "useCase.Handle",
			slog.String("details", "metadata validation"),
			slog.String("details", string(bPayloadValidationExecutionResponse)))

		if payloadValidationExecutionResponse.Failures != nil {
			return exceptions.NewMetadataValidationError(*payloadValidationExecutionResponse.Failures)
		}
	}
	return nil
}

func NewUseCase(apiPromptRoadMapConfigExecution interfaces.ApiPromptRoadMapConfigExecution, apiPromptRoadMapConfig interfaces.ApiPromptRoadMapConfig, validationApi interfaces.ApiValidation, queueAiRequester interfaces.QueueAiRequester) interfaces.UseCase {
	return &UseCase{
		apiPromptRoadMapConfig:          apiPromptRoadMapConfig,
		apiPromptRoadMapConfigExecution: apiPromptRoadMapConfigExecution,
		apiValidation:                   validationApi,
		queueAiRequester:                queueAiRequester,
	}
}
