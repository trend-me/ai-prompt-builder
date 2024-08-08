package validations

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/delivery/dtos"
)

func getError(tag string, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("'%s' is required", field)
	case "min":
		return fmt.Sprintf("'%s' should be greather in length", field)
	case "uuid":
		return fmt.Sprintf("'%s' should be an uuid", field)
	}

	return fmt.Sprintf("unmaped error for field %s", field)
}
func getField(field string) string {
	return map[string]string{
		"PromptRoadMapConfigName":        "prompt_road_map_config_name",
		"PromptRoadMapConfigExecutionId": "prompt_road_map_config_execution_id",
		"Step":                           "step",
		"OutputQueue":                    "output_queue",
		"Model":                          "model",
		"Metadata":                       "metadata"}[field]
}

func ValidateRequest(request *dtos.Request) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(request)
	if err == nil {
		return nil
	}

	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		messages = append(messages, getError(err.ActualTag(), getField(err.Field())))
	}

	return exceptions.NewValidationError(messages...)
}
