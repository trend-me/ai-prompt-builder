package steps

import "testing"
import "github.com/cucumber/godog"

func TestFeatures(t *testing.T) {
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

func aMessageShouldBeSentToTheAirequesterQueue(arg1 *godog.Table) error {
	return godog.ErrPending
}

func aMessageWithTheFollowingDataIsSentToAipromptbuilderQueue(arg1 *godog.Table) error {
	return godog.ErrPending
}

func ifTheStepIsNotEqualToTheMetadataIsSentToTheValidationAPIWithTheMetadata_validation_name(arg1 int) error {
	return godog.ErrPending
}

func noMessageShouldBeSentToTheAirequesterQueue() error {
	return godog.ErrPending
}

func noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi() error {
	return godog.ErrPending
}

func noPrompt_road_map_config_executionShouldBeUpdated() error {
	return godog.ErrPending
}

func theApplicationShouldNotRetry() error {
	return godog.ErrPending
}

func theMessageIsConsumedByTheAipromptbuilderConsumer() error {
	return godog.ErrPending
}

func thePromptRoadMapAPIReturnsAnStatusCode(arg1 int) error {
	return godog.ErrPending
}

func thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap(arg1 *godog.Table) error {
	return godog.ErrPending
}

func thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_name() error {
	return godog.ErrPending
}

func thePrompt_road_map_config_executionIsUpdatedWithTheCurrentStepOfThePrompt_road_map() error {
	return godog.ErrPending
}

func theValidationAPIReturnsAnErrorValidatingTheMetadata(arg1 *godog.Table) error {
	return godog.ErrPending
}

func theValidationAPIReturnsNoErrorsValidatingTheMetadata() error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a message should be sent to the \'ai-requester\' queue:$`, aMessageShouldBeSentToTheAirequesterQueue)
	ctx.Step(`^a message with the following data is sent to \'ai-prompt-builder\' queue:$`, aMessageWithTheFollowingDataIsSentToAipromptbuilderQueue)
	ctx.Step(`^if the step is not equal to (\d+), the metadata is sent to the validation API with the metadata_validation_name$`, ifTheStepIsNotEqualToTheMetadataIsSentToTheValidationAPIWithTheMetadata_validation_name)
	ctx.Step(`^no message should be sent to the \'ai-requester\' queue$`, noMessageShouldBeSentToTheAirequesterQueue)
	ctx.Step(`^no prompt_road_map should be fetched from the prompt-road-map-api$`, noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi)
	ctx.Step(`^no prompt_road_map_config_execution should be updated$`, noPrompt_road_map_config_executionShouldBeUpdated)
	ctx.Step(`^the application should not retry$`, theApplicationShouldNotRetry)
	ctx.Step(`^the message is consumed by the ai-prompt-builder consumer$`, theMessageIsConsumedByTheAipromptbuilderConsumer)
	ctx.Step(`^The prompt road map API returns an statusCode (\d+)$`, thePromptRoadMapAPIReturnsAnStatusCode)
	ctx.Step(`^The prompt road map API returns the following prompt road map:$`, thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap)
	ctx.Step(`^the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_name$`, thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_name)
	ctx.Step(`^the prompt_road_map_config_execution is updated with the current step of the prompt_road_map$`, thePrompt_road_map_config_executionIsUpdatedWithTheCurrentStepOfThePrompt_road_map)
	ctx.Step(`^The validation API returns an error validating the metadata:$`, theValidationAPIReturnsAnErrorValidatingTheMetadata)
	ctx.Step(`^The validation API returns no errors validating the metadata$`, theValidationAPIReturnsNoErrorsValidatingTheMetadata)
}
