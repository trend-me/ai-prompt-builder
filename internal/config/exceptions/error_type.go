package exceptions

import (
	"encoding/json"
	"strings"
)

type ErrorType struct {
	Abort     bool
	Notify    bool
	ErrorType string
	message   []string
}

func (e ErrorType) Error() string {
	return strings.Join(e.message, ";")
}

func (e ErrorType) JSON() []byte {
	j := struct {
		Abort     bool   `json:"abort"`
		Notify    bool   `json:"notify"`
		ErrorType string `json:"error_type"`
		Message   string `json:"message"`
	}{
		Abort:     e.Abort,
		Notify:    e.Notify,
		ErrorType: e.ErrorType,
		Message:   strings.Join(e.message, ";"),
	}
	b, _ := json.Marshal(j)
	return b
}
