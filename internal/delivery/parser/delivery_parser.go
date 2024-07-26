package parser

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/golang-rabbitmq-lib/rabbitmq"
)

func ParseDeliveryJSON(out interface{}, delivery amqp091.Delivery) error {
	if delivery.ContentType != rabbitmq.ContentTypeJson {
		return exceptions.NewValidationError(
			fmt.Sprintf("ContentType (%s) should be %s",
				delivery.ContentType, rabbitmq.ContentTypeJson))
	}

	return json.Unmarshal(delivery.Body, out)
}
