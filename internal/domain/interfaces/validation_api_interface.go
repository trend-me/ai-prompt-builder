package interfaces

import (
	"context"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

type ApiValidation interface {
	ExecutePayloadValidator(ctx context.Context, PayloadValidatorName string, payload []byte) (*models.PayloadValidatorExecutionResponse, error)
}
