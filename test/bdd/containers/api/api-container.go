package api_container

import (
	"bytes"
	"context"
	"fmt"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"io"
	"net/http"
)

var compose tc.ComposeStack

func Connect() error {
	c, err := tc.NewDockerCompose("docker-compose.yml")
	if err != nil {
		return err
	}
	compose = c
	return nil

}
func Disconnect() error {
	err := compose.Down(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func MockPayloadValidationsApiExecuteValidation(response string) {
	jsonString := fmt.Sprintf(`{
   "httpRequest": {
     "method": "POST",
     "path": "/payload_validations/{payload_validation_name}",
     "pathParameters": {
       "payload_validation_name": ["[A-Z0-9\\-]+"]
     }
   },
   "httpResponse": {
     "body": "%s"
   }
 })`, response)

	mockRequest(jsonString)
}

func MockPromptRoadMapConfigsApiGetPromptRoadMap(promptRoadMapConfigName string, promptRoadMapStep int, response string) {
	jsonString := fmt.Sprintf(`  {
    "httpRequest": {
      "method": "GET",
      "path": "/prompt_road_map_configs/{prompt_road_map_config_name}/prompt_road_maps/{step}",
      "pathParameters": {
        "prompt_road_map_config_name": "%s",
		"step": %d
      }
    },
    "httpResponse": {
      "body": "%s"
    }
  })`, promptRoadMapConfigName, promptRoadMapStep, response)

	mockRequest(jsonString)
}

func MockPromptRoadMapConfigExecutionsApiUpdateStepInExecution(promptRoadMapConfigExecutionId string, stepInExecution int, response string) {
	jsonString := fmt.Sprintf(`  {
    "httpRequest": {
      "method": "GET",
      "path": "/prompt_road_map_config_executions/{prompt_road_map_config_execution_id}"/step_in_execution",
      "pathParameters": {
        "prompt_road_map_config_execution_id": "%s",
      },
      "body": {
        "type": "JSON",
        "step_in_execution": %d
      }
    },
    "httpResponse": {
      "body": "%s"
    }
  })`, promptRoadMapConfigExecutionId, stepInExecution, response)

	mockRequest(jsonString)
}

func mockRequest(jsonString string) {

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
