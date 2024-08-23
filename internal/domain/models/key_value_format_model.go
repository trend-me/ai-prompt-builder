package models

import "time"

type KeyValueFormat struct {
	Id                         *string    `json:"id,omitempty"`
	PayloadValidatorName       *string  `json:"payload_validator_name,omitempty"`
	Key                        *string    `json:"key,omitempty"`
	Type                       *string    `json:"type,omitempty"`
	Match                      *string    `json:"match,omitempty"`
	Required                   *bool      `json:"required,omitempty"`
	ChildrenPayloadValidatorName *string  `json:"children_payload_validator_name,omitempty"`
	CreatedAt                  *time.Time `json:"created_at,omitempty"`
	UpdatedAt                  *time.Time `json:"updated_at,omitempty"`
}
 