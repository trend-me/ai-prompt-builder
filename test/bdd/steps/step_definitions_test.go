package steps

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/trend-me/ai-prompt-builder/internal/config/injector"
	"github.com/trend-me/ai-prompt-builder/internal/config/properties"
	"github.com/trend-me/ai-prompt-builder/test/bdd/containers"
	api_container "github.com/trend-me/ai-prompt-builder/test/bdd/containers/api"
	rabbitmq_container "github.com/trend-me/ai-prompt-builder/test/bdd/containers/rabbitmq"
	"testing"
	"time"
)
import "github.com/cucumber/godog"

var (
	t               *testing.T
	consumedMessage string
	consumer        func(ctx context.Context) (chan error, error)
)

func setup(t *testing.T) {
	err := godotenv.Load("../.bdd.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	containers.Up()
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

	consumer, err = injector.InitializeConsumer()
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

func TestFeatures(t_ *testing.T) {
	t = t_
	setup(t)
	t.Cleanup(func() {
		defer down(t)
	})

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

func aMessageWithTheFollowingDataShouldBeSentToAipromptbuilderQueue(queue string, arg1 *godog.DocString) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(queue)
	if err != nil {
		return err
	}

	assert.JSONEq(t, arg1.Content, string(content))
	return nil
}

func noMessageShouldBeSentToTheAirequesterQueue(queue string) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(queue)
	if err != nil {
		return err
	}

	assert.Nil(t, content)
	return nil
}

func noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi() error {
	return nil
}

func noPrompt_road_map_config_executionShouldBeUpdated() error {
	return nil
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh, err := consumer(ctx)
	if err != nil {
		return err
	}

	for {
		select {
		case <-errCh:
			return err
		case <-time.After(120 * time.Second):
			err = fmt.Errorf("timeout")
			return err
		}
	}
}

func theMetadataShouldBeSentToTheValidationAPIWithTheMetadata_validation_nameTEST_METADATA() error {
	return nil
}

func theMetadataShouldNotBeSentToTheValidationAPI() error {
	return nil
}

func thePromptRoadMapAPIReturnsAnStatusCode500() error {
	return api_container.MockPromptRoadMapConfigsStatus(500)
}

func thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap(step int, name string, arg1 *godog.DocString) error {
	return api_container.MockPromptRoadMapConfigsApiGetPromptRoadMap(name, step, arg1.Content)
}

func thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_config_name() error {
	_, err := api_container.GetRequests()

	return err
}

func thePrompt_road_map_config_executionIsUpdatedWithTheCurrentStepOfThePrompt_road_map() error {
	return nil
}

func theValidationAPIReturnsTheFolowingValidationResult(arg1 *godog.DocString) error {
	return api_container.MockPayloadValidationsApiExecuteValidation(arg1.Content)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		err := api_container.Reset()
		return ctx, err
	})
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
	ctx.Step(`^the prompt road map API returns an statusCode 500$`, thePromptRoadMapAPIReturnsAnStatusCode500)
	ctx.Step(`^the prompt road map API returns the following prompt road map for step \'(\d+)\' and prompt_road_map_config_name \'(.*)\':$`, thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap)
	ctx.Step(`^the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name \'(.*)\'$`, thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_config_name)
	ctx.Step(`^the prompt_road_map_config_execution step_in_execution is updated to \'(\d+)\'$`, thePrompt_road_map_config_executionIsUpdatedWithTheCurrentStepOfThePrompt_road_map)
	ctx.Step(`^the validation API returns the following validation result:$`, theValidationAPIReturnsTheFolowingValidationResult)
}
