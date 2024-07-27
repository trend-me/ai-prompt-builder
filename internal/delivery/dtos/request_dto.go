package dtos

type Request struct {
	PromptRoadMapId                string         `json:"prompt_road_map_id,omitempty" validate:"required,uuid"`
	PromptRoadMapConfigExecutionId string         `json:"prompt_road_map_config_execution_id,omitempty" validate:"required,uuid"`
	OutputQueue                    string         `json:"output_queue,omitempty" validate:"required,min=1"`
	Model                          string         `json:"model,omitempty" validate:"required,min=1"`
	Metadata                       map[string]any `json:"metadata,omitempty" validate:"required,min=1"`
}
