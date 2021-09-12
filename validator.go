package validator

import (
	"github.com/gobigbang/validator/messages"
	"github.com/gobigbang/validator/types"
)

/*
Add a validation rule to the global registry
*/
func AddValidationRule(name string, checker types.RuleCheckerFunc) {
	types.RuleRegistry[name] = checker
}

var DefaultValidatorConfig = types.ValidatorConfig{
	StructTag: "json",
	Messages:  messages.Messages,
}

func New(rules interface{}, config types.ValidatorConfig) types.Validator {
	return types.Validator{
		Rules:  rules,
		Config: config,
	}
}

func Default(rules interface{}) types.Validator {
	return New(rules, DefaultValidatorConfig)
}

/*
The function FieldRules defines the validation rules (rules ...interface{})
for a single field (field string)
*/
func FieldRules(field string, rules ...interface{}) types.FieldRules {
	return types.FieldRules{
		FieldName: field,
		Rules:     rules,
	}
}

/*
The function MapRules defines the validation for a complex type (map or struct)
recives the rules to be applied to each field
*/
func MapRules(rules ...types.FieldRules) types.MapRules {
	return types.MapRules{
		FieldRules: func() map[string]types.FieldRules {
			ret := make(map[string]types.FieldRules)
			for _, v := range rules {
				ret[v.FieldName] = v
			}
			return ret
		}(),
	}
}

func ArrayRules(rules ...types.ArrayRulesItem) types.ArrayRules {
	return types.ArrayRules{
		Rules: func() map[string]types.ArrayRulesItem {
			ret := make(map[string]types.ArrayRulesItem)
			for _, v := range rules {
				ret[v.Key] = v
			}
			return ret
		}(),
	}
}
func ArrayItemRules(key string, rules ...interface{}) types.ArrayRulesItem {
	return types.ArrayRulesItem{
		Key:   key,
		Rules: rules,
	}
}

func Conditional(check types.ConditionalRuleFunc, rules ...interface{}) types.ConditionalRule {
	return types.ConditionalRule{
		Cond:  check,
		Rules: rules,
	}
}
