package models

type Request struct {
	PromptRoadMapId                string
	PromptRoadMapConfigExecutionId string
	OutputQueue                    string
	Model                          string
	Metadata                       map[string]any
}
