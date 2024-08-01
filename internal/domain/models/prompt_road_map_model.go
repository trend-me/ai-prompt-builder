package models

import "time"

type PromptRoadMap struct {
	ResponseValidationName  *string    `json:"response_validation_id,omitempty"`
	MetadataValidationName  *string    `json:"metadata_validation_id,omitempty"`
	PromptRoadMapConfigName *string    `json:"prompt_road_map_config_name,omitempty"`
	QuestionTemplate        *string    `json:"question_template,omitempty"`
	Step                    *int32     `json:"step,omitempty"`
	CreatedAt               *time.Time `json:"created_at,omitempty"`
	UpdatedAt               *time.Time `json:"updated_at,omitempty"`
}
