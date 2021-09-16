package validator

import "github.com/gobigbang/validator/rules"

func init() {
	RegisterRule("required", rules.Required)
	RegisterRule("not_zero", rules.NotZero)
	RegisterRule("non_zero", rules.NotZero)

	// // "is" rules
	// RegisterRule("is_string", rules.IsString)
	// RegisterRule("string", rules.IsString)
	// RegisterRule("is_numeric", rules.IsNumeric)
	// RegisterRule("numeric", rules.IsNumeric)

	// //other
	// RegisterRule("in", rules.In)
	// RegisterRule("not_in", rules.NotIn)

}
