package validator

import (
	"github.com/gobigbang/validator/types"
)

// Creates a new Valiator instance (*types.Validator)
func New(rules types.RuleMap) *types.Validator {
	return types.NewValidator(rules)
}

// Registers a rule in the global registry
func RegisterRule(name string, rule types.IRule) {
	types.Registry[name] = rule
}

// Creates a new rule from the closure
func Closure(f types.RuleClosure) types.Rule {
	return types.Closure(f)
}

// Creates a single var validator RuleMap
func Var(varName string, rules ...interface{}) types.RuleMap {
	return Rules(Field("", rules...).Alias(varName))
}

// Creates a FieldRules
// The field name is the fieldpath to be validated (on child objects or keys the format is a.b.c)
// The rules can be any IRule or ruleCode (from the global registry)
func Field(fieldName string, rules ...interface{}) types.FieldRules {
	return types.FieldRules{}.Field(fieldName).Rules(rules)
}

// Creates a conditional rule
func Conditional(cond types.ConditionalRuleCond) types.ConditionalRule {
	return types.ConditionalRule{}.Condition(cond)
}

// Creates a RuleMap to be passed to the validator
func Rules(rules ...types.FieldRules) types.RuleMap {
	r := types.NewRuleMap()
	if len(rules) > 0 {
		rm := make(map[string]types.FieldRules)
		for _, v := range rules {
			rm[v.GetField()] = v
		}
		r = r.Rules(rm)
	}
	return r
}
