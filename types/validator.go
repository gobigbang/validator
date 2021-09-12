package types

import (
	"context"
	"strings"

	"github.com/gobigbang/validator/utils"
)

var RuleRegistry map[string]RuleCheckerFunc = make(map[string]RuleCheckerFunc)

// Obtains the rule funciton from the name or funcion
func getRuleFunc(r interface{}) ParsedRule {
	var ret ParsedRule

	switch r.(type) {
	case string:
		parts := utils.ParseFuncRegex(r.(string))
		if parts != nil {
			ret = ParsedRule{
				Checker: RuleRegistry[parts["RuleName"]],
				Args:    strings.Split(parts["Args"], ","),
			}
		}
	case RuleCheckerFunc:
		ret = ParsedRule{
			Checker: r.(RuleCheckerFunc),
		}
	}

	return ret
}

type (
	Validator struct {
		Rules  interface{}
		Config ValidatorConfig
	}
	// Validator configuration values
	ValidatorConfig struct {
		StructTag string `default:"json"`
		Messages  map[string]interface{}
	}
)

func (v Validator) Validate(ctx context.Context, input interface{}) error {
	return Validate(ctx, input, v.Rules, v.Config)
}
