package properties

import (
	"os"
	"strconv"
)

const (
	QueueNameAiPromptBuilder = "ai-prompt-builder"
)

func CreateQueueIfNX() bool {
	return os.Getenv("CREATE_QUEUE_IF_NX") == "true"
}

func QueueConnectionUser() string {
	return os.Getenv("QUEUE_CONNECTION_USER")
}

func QueueConnectionPort() string {
	return os.Getenv("QUEUE_CONNECTION_PORT")
}

func QueueConnectionHost() string {
	return os.Getenv("QUEUE_CONNECTION_HOST")
}

func QueueConnectionPassword() string {
	return os.Getenv("QUEUE_CONNECTION_PASSWORD")
}

func GetMaxReceiveCount() int {
	i, _ := strconv.Atoi(os.Getenv("MAX_RECEIVE_COUNT"))
	return i
}
