package exceptions

func NewValidationError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Validation Error",
		Message:   messages,
	}
}

func NewUnknownError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    false,
		ErrorType: "Unknown Error",
		Message:   messages,
	}
}

func NewMetadataValidationError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Metadata Validation Error",
		Message:   messages,
	}
}

func NewQueueError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Queue Error",
		Message:   messages,
	}
}

func NewUpdatePromptRoadMapConfigExecutionError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Update Prompt Road Map Config Execution Error",
		Message:   messages,
	}
}

func NewGetPromptRoadMapConfigError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Get Prompt Road Map Config Error",
		Message:   messages,
	}
}

func NewPromptRoadMapConfigNotFoundError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Prompt Road Map Not Found Error",
		Message:   messages,
	}
}

func NewPayloadValidatorNotFoundError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Payload Validator Not Found Error",
		Message:   messages,
	}
}

func NewPayloadValidatorError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Payload Validator Error",
		Message:   messages,
	}
}
