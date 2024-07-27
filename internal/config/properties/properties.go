package properties

import (
	"context"
	"os"
	"strconv"
)

type contextKey = string

const (
	ctxReceiveCount          contextKey = "ctxReceiveCount"
	QueueNameAiPromptBuilder            = "ai-prompt-builder"
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

func UrlApiPromptRoadMap() string {
	return os.Getenv("URL_API_PROMPT_ROAD_MAP")
}

func UrlApiValidation() string {
	return os.Getenv("URL_API_VALIDATION")
}

func QueueConnectionPassword() string {
	return os.Getenv("QUEUE_CONNECTION_PASSWORD")
}

func GetMaxReceiveCount() int {
	i, _ := strconv.Atoi(os.Getenv("MAX_RECEIVE_COUNT"))
	return i
}

func SetCtxRetryCount(ctx context.Context, receiveCount int) context.Context {
	return context.WithValue(ctx, ctxReceiveCount, receiveCount)
}

func GetCtxRetryCount(ctx context.Context) int {
	i, _ := ctx.Value(ctxReceiveCount).(int)
	return i
}
