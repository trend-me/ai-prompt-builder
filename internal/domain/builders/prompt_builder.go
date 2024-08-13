package builders

import (
	"github.com/trend-me/ai-prompt-builder/internal/config/exceptions"
	"github.com/trend-me/ai-prompt-builder/internal/domain/models"
)

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func toAnySlice(value interface{}) []any {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Slice {
		return nil
	}

	anySlice := make([]any, val.Len())
	for i := 0; i < val.Len(); i++ {
		anySlice[i] = val.Index(i).Interface()
	}

	return anySlice
}

func replaceTags(template string, metadata map[string]interface{}) (string, error) {
	re := regexp.MustCompile(`<([^>]+)>`)
	var errors string

	result := re.ReplaceAllStringFunc(template, func(tag string) string {
		key := strings.Trim(tag, "<>")      // Remove the angle brackets
		dotSplit := strings.Split(key, ".") // Split the key by dot

		var parts []string
		for _, part := range dotSplit {
			indexSplit := strings.Split(part, "[")
			if len(indexSplit) > 1 {
				indexSplit[1] = "[" + indexSplit[1]
			}

			parts = append(parts, indexSplit...)

		}

		var value any
		value = metadata

		for _, part := range parts {
			if value == nil {
				return tag // If at any point the value is nil, return the tag itself
			}
			if strings.Contains(part, "[") {
				slice := toAnySlice(value)
				if slice == nil {
					errors += fmt.Sprintf("part '%s'should be an array. tag: '%s';", part, tag)
					return tag
				}

				index, err := strconv.Atoi(part[1 : len(part)-1])
				if err != nil {
					errors += fmt.Sprintf("index '%s' is not a number. tag: '%s';", part[1:len(part)-1], tag)
					return tag
				}
				value = slice[index]
				continue
			}

			switch v := value.(type) {
			default:
				fmt.Println(value)
			case map[string]interface{}:
				value = v[part]
			}
		}

		if m, ok := value.(map[string]interface{}); ok {
			b, _ := json.Marshal(m)
			return string(b)
		}

		return fmt.Sprintf("%v", value)
	})

	if len(errors) > 0 {
		return result, exceptions.NewValidationError("error during prompt build", errors)
	}
	return result, nil
}
func BuildPrompt(request *models.Request, promptRoadMap *models.PromptRoadMap) (string, error) {
	return replaceTags(*promptRoadMap.QuestionTemplate, request.Metadata)
}
