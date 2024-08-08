Feature: Build AI prompts and forward the event to ai-requester queue

  Scenario: Successfully process a message from the queue
    Given a message with the following data is sent to 'ai-prompt-builder' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "model":"GEMINI",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    Given the prompt road map API returns the following prompt road map:
    """
    {
    "prompt_road_map_config_name":"TEST_METADATA",
    "response_validation_name":"TEST_RESPONSE",
    "metadata_validation_name":"TEST_METADATA",
    "question_template":"this is a <any.thing>. <any.array> <any.array[0]>",
    "step":2,
    "created_at":"2024-08-01T20:53:49.132Z",
    "updated_at":"2024-08-01T20:53:49.132Z"
    }
    """
    Given the validation API returns the following validation result:
    """
    {
      "failures": "",
      "errors":[]
    }
    """
    When the message is consumed by the ai-prompt-builder consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name
    And the prompt_road_map_config_execution step_in_execution is updated to '2'
    And the metadata should be sent to the validation API with the metadata_validation_name 'TEST_METADATA'
    And a message should be sent to the 'ai-requester' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "model":"GEMINI",
    "prompt":"this is a test. [1 2 3 4] 1",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    And the application should not retry

  Scenario: Successfully process a message from the queue with step 1
    Given a message with the following data is sent to 'ai-prompt-builder' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "model":"GEMINI",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    Given the prompt road map API returns the following prompt road map:
    """
    {
    "prompt_road_map_config_name":"TEST_METADATA",
    "response_validation_name":"TEST_RESPONSE",
    "metadata_validation_name":"TEST_METADATA",
    "question_template":"this is a <any.thing>. <any.array> <any.array[0]>",
    "step":1,
    "created_at":"2024-08-01T20:53:49.132Z",
    "updated_at":"2024-08-01T20:53:49.132Z"
    }
    """
    When the message is consumed by the ai-prompt-builder consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name
    And the prompt_road_map_config_execution step_in_execution is updated to '1'
    And the metadata should not be sent to the validation API
    And a message with the following data should be sent to 'ai-requester' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "model":"GEMINI",
    "prompt":"this is a test. [1 2 3 4] 1",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    And the application should not retry

  Scenario: Successfully handles a metadata validation failure
    Given a message with the following data is sent to 'ai-prompt-builder' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "model":"GEMINI",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    Given the prompt road map API returns the following prompt road map:
    """
    {
    "prompt_road_map_config_name":"TEST_METADATA",
    "response_validation_name":"TEST_RESPONSE",
    "metadata_validation_name":"TEST_METADATA",
    "question_template":"this is a <any.thing>. <any.array> <any.array[0]>",
    "step":2,
    "created_at":"2024-08-01T20:53:49.132Z",
    "updated_at":"2024-08-01T20:53:49.132Z"
    }
    """
    Given the validation API returns the following validation result:
    """
    {
      "failures": "there is something wrong...",
      "errors":[]
    }
    """
    When the message is consumed by the ai-prompt-builder consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name
    And the prompt_road_map_config_execution step_in_execution is updated to '2'
    And the metadata should be sent to the validation API with the metadata_validation_name 'TEST_METADATA'
    And no message should be sent to the 'ai-requester' queue
    And the application should not retry

  Scenario: Successfully process an error and scheduling a retry
    Given a message with the following data is sent to 'ai-prompt-builder' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "model":"GEMINI",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    Given The prompt road map API returns an statusCode 500
    When the message is consumed by the ai-prompt-builder consumer
    Then no prompt_road_map should be fetched from the prompt-road-map-api
    And no prompt_road_map_config_execution should be updated
    And the metadata should not be sent to the validation API
    And no message should be sent to the 'ai-requester' queue
    And the application should retry