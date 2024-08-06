Feature: Build AI prompts and forward the event to ai-requester queue

  Scenario: Successfully process a message from the queue
    Given a message with the following data is sent to 'ai-prompt-builder' queue:
      | prompt_road_map_config_name                   | prompt_road_map_config_id            | output_queue | model  | metadata                                       |
      | f4d28f88-5ff1-474c-842d-ad9c3ed4679b | ffcb8c0a-76aa-4490-b851-f092ed493807 | output-queue | GEMINI | {"any": { "thing":"test", "array":[1,2,3,4]} } |
    Given The prompt road map API returns the following prompt road map:
      | response_validation_id               | metadata_validation_id               | prompt_road_map_config_name | question_template                                 | step | created_at               | updated_at               |
      | c713deb9-efa2-4d5f-9675-abe0b7e0c0d4 | 7b2644c0-f5d6-43b2-a99a-cbccc2b4ab4a | BDD_PROMPT                  | this is a <any.thing>. <any.array> <any.array[0]> | 0    | 2024-08-01T20:53:49.132Z | 2024-08-01T20:53:49.132Z |
    Given The validation API returns no errors validating the metadata
    When the message is consumed by the ai-prompt-builder consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name
    And the prompt_road_map_config_execution is updated with the current step of the prompt_road_map
    And if the step is not equal to 0, the metadata is sent to the validation API with the metadata_validation_name
    And a message should be sent to the 'ai-requester' queue:
      | prompt_road_map_config_name                   | prompt_road_map_config_id            | output_queue | model  | metadata                                       | prompt                      |
      | f4d28f88-5ff1-474c-842d-ad9c3ed4679b | ffcb8c0a-76aa-4490-b851-f092ed493807 | output-queue | GEMINI | {"any": { "thing":"test", "array":[1,2,3,4]} } | this is a test. [1 2 3 4] 1|
    And the application should not retry

  Scenario: Successfully process a message from the queue when prompt road map has status != 0
    Given a message with the following data is sent to 'ai-prompt-builder' queue:
      | prompt_road_map_config_name                   | prompt_road_map_config_id            | output_queue | model  | metadata                                       |
      | f4d28f88-5ff1-474c-842d-ad9c3ed4679b | ffcb8c0a-76aa-4490-b851-f092ed493807 | output-queue | GEMINI | {"any": { "thing":"test", "array":[1,2,3,4]} } |
    Given The prompt road map API returns the following prompt road map:
      | response_validation_id               | metadata_validation_id               | prompt_road_map_config_name | question_template                                 | step | created_at               | updated_at               |
      | c713deb9-efa2-4d5f-9675-abe0b7e0c0d4 | 7b2644c0-f5d6-43b2-a99a-cbccc2b4ab4a | BDD_PROMPT                  | this is a <any.thing>. <any.array> <any.array[0]> | 4    | 2024-08-01T20:53:49.132Z | 2024-08-01T20:53:49.132Z |
    Given The validation API returns no errors validating the metadata
    When the message is consumed by the ai-prompt-builder consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name
    And the prompt_road_map_config_execution is updated with the current step of the prompt_road_map
    And if the step is not equal to 0, the metadata is sent to the validation API with the metadata_validation_name
    And a message should be sent to the 'ai-requester' queue:
      | prompt_road_map_config_name                   | prompt_road_map_config_id            | output_queue | model  | metadata                                       | prompt                      |
      | f4d28f88-5ff1-474c-842d-ad9c3ed4679b | ffcb8c0a-76aa-4490-b851-f092ed493807 | output-queue | GEMINI | {"any": { "thing":"test", "array":[1,2,3,4]} } | this is a test. [1 2 3 4] 1|
    And the application should not retry

  Scenario: Successfully handles a metadata validation failure
    Given a message with the following data is sent to 'ai-prompt-builder' queue:
      | prompt_road_map_config_name            | prompt_road_map_config_id     | output_queue  | model | metadata           |
      | "uuid"                        | "uuid"                       | "outputQueue" | GEMINI| {"any": "thing"}   |
    Given The validation API returns an error validating the metadata:
      | failures                    |
      | "error details"             |
    Given The prompt road map API returns the following prompt road map:
      | response_validation_id               | metadata_validation_id               | prompt_road_map_config_name | question_template                                 | step | created_at               | updated_at               |
      | c713deb9-efa2-4d5f-9675-abe0b7e0c0d4 | 7b2644c0-f5d6-43b2-a99a-cbccc2b4ab4a | BDD_PROMPT                  | this is a <any.thing>. <any.array> <any.array[0]> | 0    | 2024-08-01T20:53:49.132Z | 2024-08-01T20:53:49.132Z |
    When the message is consumed by the ai-prompt-builder consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name
    And the prompt_road_map_config_execution is updated with the current step of the prompt_road_map
    And if the step is not equal to 0, the metadata is sent to the validation API with the metadata_validation_name
    And no message should be sent to the 'ai-requester' queue
    And the application should not retry

  Scenario: Successfully process an error and scheduling a retry
    Given a message with the following data is sent to 'ai-prompt-builder' queue:
      | prompt_road_map_config_name            | prompt_road_map_config_id     | output_queue  | model | metadata           |
      | "uuid"                        | "uuid"                       | "outputQueue" | GEMINI| {"any": "thing"}   |
    Given The prompt road map API returns an statusCode 500
    When the message is consumed by the ai-prompt-builder consumer
    Then no prompt_road_map should be fetched from the prompt-road-map-api
    And no prompt_road_map_config_execution should be updated
    And if the step is not equal to 0, the metadata is sent to the validation API with the metadata_validation_name
    And no message should be sent to the 'ai-requester' queue
    And the application should not retry