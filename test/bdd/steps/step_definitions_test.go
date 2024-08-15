package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/trend-me/ai-prompt-builder/internal/domain/interfaces"

	"github.com/cucumber/godog"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/trend-me/ai-prompt-builder/internal/config/injector"
	"github.com/trend-me/ai-prompt-builder/internal/config/properties"
	"github.com/trend-me/ai-prompt-builder/test/bdd/containers"
	rabbitmq_container "github.com/trend-me/ai-prompt-builder/test/bdd/containers/rabbitmq"
	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/params"
	"github.com/vitorsalgado/mocha/v3/reply"
)

var (
	t                                                 *testing.T
	consumedMessage                                   string
	consumer                                          interfaces.QueueAiPromptBuilderConsumer
	m                                                 *mocha.Mocha
	scopePromptRoadMapConfigsApiGetPromptRoadMap      *mocha.Scoped
	scopePromptRoadMapConfigExecutionsApiUpdateStep   *mocha.Scoped
	scopePayloadValidationApiExecute                  *mocha.Scoped
	requestPayloadValidationApiExecute                *http.Request
	requestPromptRoadMapConfigExecutionsApiUpdateStep *http.Request
	requestPromptRoadMapConfigsApiGetPromptRoadMap    *http.Request
)

func jsonEqual(a, b string) bool {
	var j1, j2 map[string]interface{}

	if err := json.Unmarshal([]byte(a), &j1); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &j2); err != nil {
		return false
	}

	return reflect.DeepEqual(j1, j2)
}

func setup(t *testing.T) {
	m = mocha.New(t)
	m.Start()
	_ = os.Setenv("URL_API_PROMPT_ROAD_MAP_CONFIG", m.URL()+"/prompt_road_map_configs")
	_ = os.Setenv("URL_API_PROMPT_ROAD_MAP_CONFIG_EXECUTION", m.URL()+"/prompt_road_map_config_executions")
	_ = os.Setenv("URL_API_VALIDATION", m.URL()+"/payload_validations")
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

	err = m.Close()
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

	if !jsonEqual(arg1.Content, string(content)) {
		return fmt.Errorf("message sent to queue '%s' is not equal to the expected message: %s. Got: %s",
			queue, arg1.Content, string(content))
	}
	return nil
}

func noMessageShouldBeSentToTheAirequesterQueue(queue string) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(queue)
	if err != nil {
		return err
	}

	if !assert.Nil(t, content){
		return fmt.Errorf("a message was sent to queue '%s'. Got: %s",
			queue, string(content))
	}
	return nil
}

func noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi() error {
	if scopePromptRoadMapConfigsApiGetPromptRoadMap.Called() {
		return fmt.Errorf("prompt road map was fetched")
	}

	return nil
}

func noPrompt_road_map_config_executionShouldBeUpdated() error {
	if scopePromptRoadMapConfigExecutionsApiUpdateStep.Called() {
		return fmt.Errorf("prompt road map config execution step was updated")
	}

	return nil
}

func theApplicationShouldNotRetry() error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(properties.QueueNameAiPromptBuilder)
	if err != nil {
		return err
	}
	assert.Nil(t, content)
	return nil
}

func theApplicationShouldRetry() error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(properties.QueueNameAiPromptBuilder)
	if err != nil {
		return err
	}
	if !jsonEqual(string(content), consumedMessage) {
		return fmt.Errorf("message sent to queue '%s' is not equal to the expected message: %s. Got: %s",
			properties.QueueNameAiPromptBuilder, consumedMessage, string(content))
	}
	return nil
}

func theMessageIsConsumedByTheAipromptbuilderConsumer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh, err := consumer.Consume(ctx)
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

func theMetadataShouldBeSentToTheValidationAPIWithTheMetadata_validation_nameTEST_METADATA(name string) error {
	if !scopePayloadValidationApiExecute.Called() || requestPayloadValidationApiExecute == nil {
		return fmt.Errorf("metadata was not sent to the validation API")
	}

	split := strings.Split(requestPayloadValidationApiExecute.URL.Path, "/")
	if split[len(split)-1] != name {
		return fmt.Errorf("metadata was not sent to the validation API with the correct metadata_validation_name. I was %s",
			requestPayloadValidationApiExecute.URL.Path)
	}

	return nil
}

func theMetadataShouldNotBeSentToTheValidationAPI() error {
	if scopePayloadValidationApiExecute.Called() {
		return fmt.Errorf("metadata was sent to the validation API")
	}

	return nil
}

func thePromptRoadMapAPIReturnsAnStatusCode500() error {
	scopePromptRoadMapConfigsApiGetPromptRoadMap = m.AddMocks(mocha.Get(expect.Func(func(v any, a expect.Args) (bool, error) {
		return strings.Contains(a.RequestInfo.Request.URL.Path, "/prompt_road_map_configs"), nil
	})).
		Reply(reply.InternalServerError().BodyString(`{"error": "Internal Server Error"}`)))
	return nil
}

func thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap(step int, name string, arg1 *godog.DocString) error {
	scopePromptRoadMapConfigsApiGetPromptRoadMap = m.AddMocks(mocha.
		Get(expect.URLPath(fmt.Sprintf("/prompt_road_map_configs/%s/prompt_road_maps/%d", name, step))).ReplyFunction(func(request *http.Request, r reply.M, p params.P) (*reply.Response, error) {
		requestPromptRoadMapConfigsApiGetPromptRoadMap = request
		return &reply.Response{
			Status: http.StatusOK,
			Body:   io.NopCloser(strings.NewReader(arg1.Content)),
		}, nil
	}))

	return nil
}

func thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_config_name(name string, step int) error {
	if !scopePromptRoadMapConfigsApiGetPromptRoadMap.Called() {
		return fmt.Errorf("prompt road map was not fetched")
	}

	if !strings.Contains(requestPromptRoadMapConfigsApiGetPromptRoadMap.URL.Path, fmt.Sprintf("/prompt_road_map_configs/%s/prompt_road_maps/%d", name, step)) {
		return fmt.Errorf("prompt_road_map_config fetched with '%s'. Requierd prompt_road_map_config_name: '%s' and step: '%d'",
			requestPromptRoadMapConfigsApiGetPromptRoadMap.URL.Path, name, step)
	}

	return nil
}

func thePrompt_road_map_config_executionIsUpdatedWithTheCurrentStepOfThePrompt_road_map(step int) error {

	if !scopePromptRoadMapConfigExecutionsApiUpdateStep.Called() || requestPromptRoadMapConfigExecutionsApiUpdateStep == nil {
		return fmt.Errorf("prompt_road_map_config_executions step was not updated")
	}

	body := make(map[string]interface{})
	decoder := json.NewDecoder(requestPromptRoadMapConfigExecutionsApiUpdateStep.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return fmt.Errorf("")
	}
	updatedStep, ok := body["step_in_execution"]
	if !ok {
		return fmt.Errorf("prompt_road_map_config_executions step was not updated with:  '%d'", step)
	}

	if updatedStepInt, ok := updatedStep.(float64); ok {
		if int(updatedStepInt) != step {
			return fmt.Errorf("prompt_road_map_config_executions step was updated with:  '%d'", int(updatedStepInt))
		}
	} else if updatedStep != step {
		return fmt.Errorf("prompt_road_map_config_executions step was updated with:  '%s'", updatedStep)
	}

	return nil
}

func theValidationAPIReturnsTheFolowingValidationResult(name string, arg1 *godog.DocString) error {
	scopePayloadValidationApiExecute = m.AddMocks(mocha.
		Post(expect.URLPath(fmt.Sprintf("/payload_validations/%s", name))).
		ReplyFunction(func(request *http.Request, r reply.M, p params.P) (*reply.Response, error) {
			requestPayloadValidationApiExecute = request
			return &reply.Response{
				Status: http.StatusOK,
				Body:   io.NopCloser(strings.NewReader(arg1.Content)),
			}, nil
		}))

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		scopePromptRoadMapConfigsApiGetPromptRoadMap = nil
		scopePromptRoadMapConfigExecutionsApiUpdateStep = nil
		scopePayloadValidationApiExecute = nil
		requestPayloadValidationApiExecute = nil
		requestPromptRoadMapConfigExecutionsApiUpdateStep = nil
		requestPromptRoadMapConfigsApiGetPromptRoadMap = nil
		consumedMessage = ""

		scopePromptRoadMapConfigExecutionsApiUpdateStep = m.AddMocks(mocha.
			Patch(expect.Func(func(v any, a expect.Args) (bool, error) {
				if !strings.Contains(a.RequestInfo.Request.URL.Path, "/prompt_road_map_config_executions") {
					return false, nil
				}

				return true, nil
			})).
			ReplyFunction(func(request *http.Request, r reply.M, p params.P) (*reply.Response, error) {
				requestPromptRoadMapConfigExecutionsApiUpdateStep = request
				return &reply.Response{
					Status: http.StatusOK,
				}, nil
			}))
		return ctx, nil
	})

	ctx.Step(`^a message with the following data is sent to \'(.*)\' queue:$`, aMessageWithTheFollowingDataIsSentToAipromptbuilderQueue)
	ctx.Step(`^a message with the following data should be sent to \'(.*)\' queue:$`, aMessageWithTheFollowingDataShouldBeSentToAipromptbuilderQueue)
	ctx.Step(`^no message should be sent to the \'(.*)\' queue:$`, noMessageShouldBeSentToTheAirequesterQueue)
	ctx.Step(`^no prompt_road_map should be fetched from the prompt-road-map-api$`, noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi)
	ctx.Step(`^no prompt_road_map_config_execution should be updated$`, noPrompt_road_map_config_executionShouldBeUpdated)
	ctx.Step(`^the application should not retry$`, theApplicationShouldNotRetry)
	ctx.Step(`^the application should retry$`, theApplicationShouldRetry)
	ctx.Step(`^the message is consumed by the ai-prompt-builder consumer$`, theMessageIsConsumedByTheAipromptbuilderConsumer)
	ctx.Step(`^the metadata should be sent to the validation API with the metadata_validation_name \'(.*)\'$`, theMetadataShouldBeSentToTheValidationAPIWithTheMetadata_validation_nameTEST_METADATA)
	ctx.Step(`^the metadata should not be sent to the validation API$`, theMetadataShouldNotBeSentToTheValidationAPI)
	ctx.Step(`^the prompt road map API returns an statusCode 500$`, thePromptRoadMapAPIReturnsAnStatusCode500)
	ctx.Step(`^the prompt road map API returns the following prompt road map for step \'(\d+)\' and prompt_road_map_config_name \'(.*)\':$`, thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap)
	ctx.Step(`^the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name \'(.*)\' and step \'(\d+)\'$`, thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_config_name)
	ctx.Step(`^the prompt_road_map_config_execution step_in_execution is updated to \'(\d+)\'$`, thePrompt_road_map_config_executionIsUpdatedWithTheCurrentStepOfThePrompt_road_map)
	ctx.Step(`^the validation API returns the following validation result for payload_validation \'(.*)\':$`, theValidationAPIReturnsTheFolowingValidationResult)
}
