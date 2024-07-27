package parsers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/config/properties"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
)

func ParseDeliveryJSON(ctx context.Context, out interface{}, delivery amqp091.Delivery) (context.Context, error) {
	retryCount, ok := delivery.Headers["x-retry-count"].(int)
	if !ok {
		retryCount = 0
	}

	ctx = properties.SetCtxRetryCount(ctx, retryCount)

	if delivery.ContentType != rabbitmq.ContentTypeJson {
		return ctx, exceptions.NewValidationError(
			fmt.Sprintf("ContentType (%s) should be %s",
				delivery.ContentType, rabbitmq.ContentTypeJson))
	}

	return ctx, json.Unmarshal(delivery.Body, out)
}
