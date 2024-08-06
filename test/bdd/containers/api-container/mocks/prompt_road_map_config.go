package mocks

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func MockPromptRoadMapConfigsApiGetPromptRoadMap(payloadValidationName, response string) {
	jsonString := fmt.Sprintf(`  {
    "httpRequest": {
      "method": "GET",
      "path": "/prompt_road_map_configs/{prompt_road_map_config_name}/prompt_road_maps/{step}",
      "pathParameters": {
        "prompt_road_map_config_name": ["[A-Z0-9\\-]+"],
        "step":["[0-9\\-]+"]
      }
    },
    "httpResponse": {
      "body": {
        "response_validation_name": "test",
        "metadata_validation_name": "valid",
        "prompt_road_map_config_name": "test",
        "question_template": "",
        "step": 1,
        "created_at": "2024-08-06T01:35:49.004Z",
        "updated_at": "2024-08-06T01:35:49.004Z"
      }
    }
  })`, payloadValidationName, response)

	req, err := http.NewRequest("POST", "https://localhost:1080/mockserver/expectation", bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(string(body))
}

/*
[
  {
    "httpRequest": {
      "method": "GET",
      "path": "/prompt_road_map_configs/{prompt_road_map_config_name}/prompt_road_maps/{step}",
      "pathParameters": {
        "prompt_road_map_config_name": [
          "valid_payload_validation"
        ],
        "step": [
          1
        ]
      }
    },
    "httpResponse": {
      "body": {
        "response_validation_name": "test",
        "metadata_validation_name": "valid",
        "prompt_road_map_config_name": "test",
        "question_template": "",
        "step": 1,
        "created_at": "2024-08-06T01:35:49.004Z",
        "updated_at": "2024-08-06T01:35:49.004Z"
      }
    }
  },
  {
    "httpRequest": {
      "method": "GET",
      "path": "/prompt_road_map_configs/{prompt_road_map_config_name}/prompt_road_maps/{step}",
      "pathParameters": {
        "prompt_road_map_config_name": [
          "invalid_payload_validation"
        ],
        "step": [
          1
        ]
      }
    },
    "httpResponse": {
      "body": {
        "response_validation_name": "test",
        "metadata_validation_name": "invalid",
        "prompt_road_map_config_name": "test",
        "question_template": "",
        "step": 1,
        "created_at": "2024-08-06T01:35:49.004Z",
        "updated_at": "2024-08-06T01:35:49.004Z"
      }
    }
  }
]
*/
