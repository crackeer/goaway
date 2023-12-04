package gateway

import (
	"strings"

	"github.com/tidwall/gjson"
)

func isKid(value string) bool {
	return strings.HasPrefix(value, "#:")
}

func getKid(value string) string {
	return value[2:]
}

func extractValue(input []byte, value string) interface{} {
	if !isKid(value) {
		return value
	}
	return gjson.GetBytes(input, getKid(value)).Value()
}

func jsonValue(dest interface{}, input []byte) interface{} {
	if value, ok := dest.(string); ok {
		return extractValue(input, value)
	}

	if mapData, ok := dest.(map[string]interface{}); ok {
		retData := make(map[string]interface{})
		for key, value := range mapData {
			retData[key] = jsonValue(value, input)
		}
		return retData
	}

	if list, ok := dest.([]interface{}); ok {
		retData := []interface{}{}
		for _, item := range list {
			retData = append(retData, jsonValue(item, input))
		}
	}
	return dest
}
