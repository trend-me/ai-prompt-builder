package exceptions

func NewValidationError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Validation Error",
		message:   messages,
	}
}

func NewUnknownError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    false,
		ErrorType: "Unknown Error",
		message:   messages,
	}
}
