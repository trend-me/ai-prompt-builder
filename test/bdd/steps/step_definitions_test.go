package steps

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/trend-me/ai-prompt-builder/internal/config/injector"
	"github.com/trend-me/ai-prompt-builder/internal/config/properties"
	"github.com/trend-me/ai-prompt-builder/test/bdd/containers"
	rabbitmq_container "github.com/trend-me/ai-prompt-builder/test/bdd/containers/rabbitmq"
	"testing"
	"time"
)
import "github.com/cucumber/godog"

var (
	consumedMessage string
	consumer        func()
)

func setup(t *testing.T) {
	containers.Up()
	var err error
	for range 10 {
		time.Sleep(10 * time.Second)
		fmt.Println("Waiting for rabbitmq to start")
		err = rabbitmq_container.Connect()
		if err == nil {
			break
		}
		fmt.Println(err.Error())
	}

	if err != nil {
		t.Fatal(err.Error())
	}

	consumer, err = injector.InitializeConsumer(context.Background())
	if err != nil {
		t.Fatal(err.Error())
	}
}

func down(t *testing.T) {
	err := rabbitmq_container.Disconnect()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = containers.Down()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFeatures(t *testing.T) {
	setup(t)
	defer down(t)
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func aMessageWithTheFollowingDataIsSentToAipromptbuilderQueue(queue string, arg1 *godog.DocString) error {
	if queue == properties.QueueNameAiPromptBuilder {
		consumedMessage = arg1.Content
	}
	return rabbitmq_container.PostMessageToQueue(queue, []byte(arg1.Content))
}

func aMessageWithTheFollowingDataShouldBeSentToAipromptbuilderQueue(t *testing.T, queue string, arg1 *godog.DocString) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(queue)
	if err != nil {
		return err
	}

	assert.JSONEq(t, arg1.Content, string(content))
	return nil
}

func noMessageShouldBeSentToTheAirequesterQueue(t *testing.T, queue string) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(queue)
	if err != nil {
		return err
	}

	assert.Nil(t, content)
	return nil
}

func noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi() error {
	return godog.ErrPending
}

func noPrompt_road_map_config_executionShouldBeUpdated() error {
	return godog.ErrPending
}

func theApplicationShouldNotRetry(t *testing.T) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(properties.QueueNameAiPromptBuilder)
	if err != nil {
		return err
	}
	assert.Nil(t, content)
	return nil
}

func theApplicationShouldRetry(t *testing.T) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(properties.QueueNameAiPromptBuilder)
	if err != nil {
		return err
	}
	assert.JSONEq(t, string(content), consumedMessage)
	return nil
}

func theMessageIsConsumedByTheAipromptbuilderConsumer() error {
	return consumer()
}

func theMetadataShouldBeSentToTheValidationAPIWithTheMetadata_validation_nameTEST_METADATA() error {
	return godog.ErrPending
}

func theMetadataShouldNotBeSentToTheValidationAPI() error {
	return godog.ErrPending
}

func thePromptRoadMapAPIReturnsAnStatusCode(arg1 int) error {
	return godog.ErrPending
}

func thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap(arg1 *godog.DocString) error {
	return godog.ErrPending
}

func thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_config_name() error {
	return godog.ErrPending
}

func thePrompt_road_map_config_executionIsUpdatedWithTheCurrentStepOfThePrompt_road_map() error {
	return godog.ErrPending
}

func theValidationAPIReturnsTheFolowingValidationResult(arg1 *godog.DocString) error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a message with the following data is sent to \'(.*)\' queue:$`, aMessageWithTheFollowingDataIsSentToAipromptbuilderQueue)
	ctx.Step(`^a message with the following data should be sent to \'(.*)\' queue:$`, aMessageWithTheFollowingDataShouldBeSentToAipromptbuilderQueue)
	ctx.Step(`^no message should be sent to the \'(.*)\' queue$`, noMessageShouldBeSentToTheAirequesterQueue)
	ctx.Step(`^no prompt_road_map should be fetched from the prompt-road-map-api$`, noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi)
	ctx.Step(`^no prompt_road_map_config_execution should be updated$`, noPrompt_road_map_config_executionShouldBeUpdated)
	ctx.Step(`^the application should not retry$`, theApplicationShouldNotRetry)
	ctx.Step(`^the application should retry$`, theApplicationShouldRetry)
	ctx.Step(`^the message is consumed by the ai-prompt-builder consumer$`, theMessageIsConsumedByTheAipromptbuilderConsumer)
	ctx.Step(`^the metadata should be sent to the validation API with the metadata_validation_name \'(.*)\'$`, theMetadataShouldBeSentToTheValidationAPIWithTheMetadata_validation_nameTEST_METADATA)
	ctx.Step(`^the metadata should not be sent to the validation API$`, theMetadataShouldNotBeSentToTheValidationAPI)
	ctx.Step(`^The prompt road map API returns an statusCode 500$`, thePromptRoadMapAPIReturnsAnStatusCode)
	ctx.Step(`^the prompt road map API returns the following prompt road map:$`, thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap)
	ctx.Step(`^the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name \'(.*)\'$`, thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_config_name)
	ctx.Step(`^the prompt_road_map_config_execution step_in_execution is updated to '(\d+)'$`, thePrompt_road_map_config_executionIsUpdatedWithTheCurrentStepOfThePrompt_road_map)
	ctx.Step(`^the validation API returns the following validation result:$`, theValidationAPIReturnsTheFolowingValidationResult)
}