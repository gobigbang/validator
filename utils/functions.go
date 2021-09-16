package utils

import "encoding/json"

func JsonPrettyPrint(input interface{}) string {
	str, _ := json.MarshalIndent(input, "", "   ")
	return string(str)
}

func Flattern(m [][]interface{}) []interface{} {
	return m[0][:cap(m[0])]
}

func ParseMessages(messages map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for i, v := range messages {
		switch v.(type) {
		case string:
			m[i] = v.(string)
		default:
			subMessages := ParseMessages(v.(map[string]interface{}))
			for si, sv := range subMessages {
				m[i+"."+si] = sv
			}
		}
	}
	return m
}
