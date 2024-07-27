package models

type PayloadValidatorExecutionResponse struct {
	// All error messages separated by ';'
	Failures *string                                   `json:"failures,omitempty"`
	Errors   []PayloadValidatorExecutionResponseErrors `json:"errors,omitempty"`
}

type PayloadValidatorExecutionResponseErrors struct {
	// the field that failed
	Field *string `json:"field,omitempty"`
	// the failure reason
	Fail           *string         `json:"fail,omitempty"`
	KeyValueFormat *KeyValueFormat `json:"key_value_format,omitempty"`
}
