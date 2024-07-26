package models

type Request struct {
	PromptRoadMapId string
	OutputQueue     string
	Model           string
	Metadata        map[string]any
}
