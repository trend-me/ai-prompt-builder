package exceptions

import (
	"encoding/json"
	"strings"
)

type ErrorType struct {
	Abort     bool     `json:"abort"`
	Notify    bool     `json:"notify"`
	ErrorType string   `json:"error_type"`
	Message   []string `json:"message"`
}

func (e ErrorType) Error() string {
	return strings.Join(e.Message, ";")
}

func (e ErrorType) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}
