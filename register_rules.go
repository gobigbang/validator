package validator

import "github.com/gobigbang/validator/rules"

func init() {
	// required rules
	AddValidationRule("required", rules.Required)

	// length rules
	AddValidationRule("min_length", rules.MinLength)
	AddValidationRule("max_length", rules.MaxLength)

	// is rules
	AddValidationRule("is_array", rules.IsArray)
	AddValidationRule("array", rules.IsArray)
	AddValidationRule("is_map", rules.IsMap)
	AddValidationRule("map", rules.IsMap)
	AddValidationRule("is_string", rules.IsString)
	AddValidationRule("string", rules.IsString)
	AddValidationRule("is_integer", rules.IsInteger)
	AddValidationRule("integer", rules.IsInteger)
	AddValidationRule("is_struct", rules.IsStruct)
	AddValidationRule("struct", rules.IsStruct)
}
