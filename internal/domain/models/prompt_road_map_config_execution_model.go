package models

type PromptRoadMapConfigExecution struct {
	Id                    *string                `json:"id,omitempty"`
	TotalSteps            *int32                 `json:"total_steps,omitempty"`
	StepInExecution       *int32                 `json:"step_in_execution,omitempty"`
	PromptRoadMapConfigId *string                `json:"prompt_road_map_config_id,omitempty"`
	Metadata              map[string]interface{} `json:"metadata,omitempty"`
}
