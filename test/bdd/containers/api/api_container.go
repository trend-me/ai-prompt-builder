package api_container

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func escapeJSONString(input string) string {
	return strings.ReplaceAll(input, `"`, `\"`)
}

type Request struct {
	Secure         bool   `json:"secure"`
	KeepAlive      bool   `json:"keepAlive"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	PathParameters []struct {
		Name   string   `json:"name"`
		Values []string `json:"values"`
	} `json:"pathParameters"`
	QueryStringParameters []struct {
		Name   string   `json:"name"`
		Values []string `json:"values"`
	} `json:"queryStringParameters"`
	Body struct {
		Not         bool   `json:"not"`
		Type        string `json:"type"`
		Base64Bytes string `json:"base64Bytes"`
		ContentType string `json:"contentType"`
	} `json:"body"`
	Headers []struct {
		Name   string   `json:"name"`
		Values []string `json:"values"`
	} `json:"headers"`
	Cookies []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"cookies"`
	SocketAddress struct {
		Host   string `json:"host"`
		Port   int    `json:"port"`
		Scheme string `json:"scheme"`
	} `json:"socketAddress"`
	Protocol string `json:"protocol"`
}

func Reset() error {
	_, err := http.NewRequest(http.MethodPut, "http://localhost:1080/mockserver/clear", nil)
	return err
}

func GetRequests() (*[]Request, error) {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", "http://localhost:1080/mockserver/retrieve", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	requests := &[]map[string]interface{}{}
	err = json.NewDecoder(resp.Body).Decode(requests)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func MockPayloadValidationsApiExecuteValidation(response string) error {
	jsonString := fmt.Sprintf(`{
   "httpRequest": {
     "method": "POST",
     "path": "/payload_validations/{payload_validation_name}",
     "pathParameters": {
        	"payload_validation_name":["[A-Z0-9\\-]+"]
     }
    },
   "httpResponse": {
     "body": "%s"
   }
 }`, escapeJSONString(response))

	return mockRequest(jsonString)
}

func MockPromptRoadMapConfigsApiGetPromptRoadMap(promptRoadMapConfigName string, promptRoadMapStep int, response string) error {
	jsonString := fmt.Sprintf(`{
  "httpRequest": {
    "method": "GET",
    "path": "/prompt_road_map_configs/{prompt_road_map_config_name}/prompt_road_maps/{step}",
    "pathParameters": {
      "prompt_road_map_config_name": ["%s"],
      "step": ["%d"]
    } },
  "httpResponse": {
    "body": "%s"
  }
}`, promptRoadMapConfigName, promptRoadMapStep, escapeJSONString(response))

	return mockRequest(jsonString)
}

func MockPromptRoadMapConfigsStatus(status int) error {
	jsonString := fmt.Sprintf(`  {
    "httpRequest": {
      "method": "GET",
      "path": "/prompt_road_map_configs/{prompt_road_map_config_name}/prompt_road_maps/{step}",
      "pathParameters": {
        	"prompt_road_map_config_name":["[A-Z0-9\\-]+"],
			"step":["[0-9]+"]
		}
    },
    "httpResponse": {
      "statusCode":%d
    }
  }`, status)

	return mockRequest(jsonString)
}

func MockPromptRoadMapConfigExecutionsApiUpdateStepInExecution(promptRoadMapConfigExecutionId string, stepInExecution int, response string) error {
	jsonString := fmt.Sprintf(`  {
    "httpRequest": {
      "method": "GET",
      "path": "/prompt_road_map_config_executions/{prompt_road_map_config_execution_id}"/step_in_execution",
      "pathParameters": {
        	"prompt_road_map_config_execution_id":["%s"],
		},
      "body": {
        "type": "JSON",
        "step_in_execution": %d
      }
    },
    "httpResponse": {
      "body":%s
    }
  }`, promptRoadMapConfigExecutionId, stepInExecution, escapeJSONString(response))

	return mockRequest(jsonString)
}

func mockRequest(jsonString string) error {

	req, err := http.NewRequest("PUT", "http://localhost:1080/mockserver/expectation", bytes.NewBuffer([]byte(jsonString)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}
	fmt.Println(string(body))
	return nil
}
