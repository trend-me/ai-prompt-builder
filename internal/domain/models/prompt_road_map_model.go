package models

import "time"

type PromptRoadMap struct {
	ResponseValidationId *string    `json:"response_validation_id,omitempty"`
	MetadataValidationId *string    `json:"metadata_validation_id,omitempty"`
	ResearchConfigId     *string    `json:"research_config_id,omitempty"`
	QuestionTemplate     *string    `json:"question_template,omitempty"`
	Step                 *int32     `json:"step,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}
