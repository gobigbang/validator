package messages

import "fmt"

var Messages = map[string]interface{}{
	"no_validation_input": "Nothing to validate",
	"key_not_allowed":     "The field {Field} is not allowed here.",
	"invalid_rule":        "The rule {Rule} is invalid.",
	"required":            "The field {Field} is required.",
	"is_array":            "The field {Field} must be an array.",
	"is_string":           "The field {Field} must be a string.",
	"is_numeric":          "The field {Field} must be a number.",
	"is_integer":          "The field {Field} must be an integer.",
	"is_map":              "The field {Field} must be a key:value map.",
	"is_struct":           "The field {Field} must be an object.",
	"max_length": map[string]string{
		"string": "The field {Field} may not be greater than {Arg0} characters long.",
		"array":  "The field {Field} may not have more than {Arg0} items.",
		"map":    "The field {Field} may not have more than {Arg0} items.",
	},
	"min_length": map[string]string{
		"string": "The field {Field} must be at least {Arg0} characters long.",
		"array":  "The field {Field} must have at least {Arg0} items.",
		"map":    "The field {Field} must have at least {Arg0} items.",
	},
}

func GetMessage(messages map[string]interface{}, code string) string {
	if m, ok := messages[code]; ok {
		return fmt.Sprint(m)
	} else {
		return code
	}
}

func GetTypedMessage(messages map[string]interface{}, code string, t string) string {
	if m, ok := messages[code]; ok {
		switch m.(type) {
		case string:
			return m.(string)
		case map[string]string:
			typed := m.(map[string]string)
			if tm, ok := typed[t]; ok {
				return tm
			}
			if tm, ok := typed["default"]; ok {
				return tm
			}
			return code
		default:
			return fmt.Sprint(m)
		}
	} else {
		return code
	}
}
