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

func NewMetadataValidationError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "Metadata Validation Error",
		message:   messages,
	}
}

func NewQueueError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Queue Error",
		message:   messages,
	}
}

func NewUpdatePromptRoadMapConfigExecutionError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Update PromptRoadMapConfigExecution Error",
		message:   messages,
	}
}

func NewGetPromptRoadMapError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     false,
		Notify:    true,
		ErrorType: "Get PromptRoadMap Error",
		message:   messages,
	}
}

func NewPromptRoadMapConfigNotFoundError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "PromptRoadMap NotFound Error",
		message:   messages,
	}
}

func NewPayloadValidatorNotFoundError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "PayloadValidator NotFound Error",
		message:   messages,
	}
}

func NewPayloadValidatorError(messages ...string) ErrorType {
	return ErrorType{
		Abort:     true,
		Notify:    true,
		ErrorType: "PayloadValidator Error",
		message:   messages,
	}
}
