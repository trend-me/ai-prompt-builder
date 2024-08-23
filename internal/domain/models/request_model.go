package models

import "github.com/trend-me/ai-prompt-builder/internal/config/exceptions"

type Request struct {
	PromptRoadMapConfigName        string
	PromptRoadMapStep              int
	PromptRoadMapConfigExecutionId string
	OutputQueue                    string
	Error                          *exceptions.ErrorType
	Model                          string
	Metadata                       map[string]any
}
