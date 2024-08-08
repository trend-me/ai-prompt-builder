package dtos

type Request struct {
	PromptRoadMapConfigName        string         `json:"prompt_road_map_config_name,omitempty" validate:"required,min=1"`
	PromptRoadMapStep              int            `json:"prompt_road_map_step,omitempty" validate:"required"`
	PromptRoadMapConfigExecutionId string         `json:"prompt_road_map_config_execution_id,omitempty" validate:"required,uuid"`
	OutputQueue                    string         `json:"output_queue,omitempty" validate:"required,min=1"`
	Model                          string         `json:"model,omitempty" validate:"required,min=1"`
	Metadata                       map[string]any `json:"metadata,omitempty" validate:"required,min=1"`
}
