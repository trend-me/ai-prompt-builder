package mocks

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

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

//[
//  {
//    "method": "POST",
//    "httpRequest": {
//      "path": "/payload_validations/{payload_validation_name}",
//      "pathParameters": {
//        "payload_validation_name": [
//          "invalid"
//        ]
//      }
//    },
//    "httpResponse": {
//      "body": {
//        "failures": "key 'testing' is required;",
//        "errors": [
//          {
//            "field": "testing",
//            "fail": "key 'testing' is required;",
//            "key_value_format": {
//              "payload_validator_name": "test",
//              "key": "testing",
//              "type": "number",
//              "match": null,
//              "required": true
//            }
//          }
//        ]
//      }
//    }
//  }
//]
