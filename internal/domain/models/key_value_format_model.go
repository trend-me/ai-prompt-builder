package models

import "time"

type KeyValueFormat struct {
	Id                         *string    `json:"id,omitempty"`
	PayloadValidatorId         *string    `json:"payload_validator_id,omitempty"`
	Key                        *string    `json:"key,omitempty"`
	Type                       *string    `json:"type,omitempty"`
	Match                      *string    `json:"match,omitempty"`
	Required                   *bool      `json:"required,omitempty"`
	ChildrenPayloadValidatorId *string    `json:"children_payload_validator_id,omitempty"`
	CreatedAt                  *time.Time `json:"created_at,omitempty"`
	UpdatedAt                  *time.Time `json:"updated_at,omitempty"`
}
