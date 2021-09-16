package types

import (
	"context"

	"github.com/gobigbang/validator/utils"
)

type (
	// Global string=> rule container
	RulesRegistry map[string]IRule

	// Request passed to the validation method
	ValidateRequest struct {
		// The field value
		Value interface{}
		// The value descriptor
		ValueDescriptor utils.ValueDescriptor
		// The root object to be validated
		RootValue interface{}
		// The validator intance
		Validator *Validator
		// The validation context
		Context context.Context
	}

	// Interface that must be implemented for the validatable objects
	IValidation interface {
		Execute(request ValidateRequest) (passes bool, e error)
	}

	// Rule interface
	IRule interface {
		Execute(request ValidateRequest) (passes bool, e error)
		GetMessage() string
		GetCode() string
		GetClosure() RuleClosure
	}

	// Rule closure
	RuleClosure func(params ValidateRequest) (passes bool, e error)

	// Rule object
	Rule struct {
		// The rule code (used for errors)
		code string
		// The error message
		message string
		// The validation method
		closure RuleClosure
	}

	RuleWithArgs struct {
		Rule
		args []interface{}
	}

	RuleWithArgsCtr func(params ...interface{}) RuleWithArgs

	Rules []Rule

	AcceptableRules []interface{}

	FieldRules struct {
		field            string
		rules            AcceptableRules
		stopOnFirstError bool
		alias            string
		messages         map[string]string
	}

	ConditionParams struct {
		Validator *Validator
	}
	ConditionalRuleCond func(ctx context.Context, value interface{}, params ConditionParams) bool

	ConditionalRule struct {
		condition ConditionalRuleCond
		failRules AcceptableRules
		rules     AcceptableRules
	}
)

func (f RuleClosure) Execute(request ValidateRequest) (passes bool, e error) {
	passes, e = f(request)
	return
}

func Closure(f RuleClosure) Rule {
	return NewRule(f)
}

func NewRule(f RuleClosure) Rule {
	return Rule{
		closure: f,
	}
}

func ClosureWithArgs(f RuleClosure, params ...interface{}) RuleWithArgs {
	return RuleWithArgs{
		Rule: Closure(f),
		args: params,
	}
}

// Rule methods

func (r Rule) Message(m string) Rule {
	r.message = m
	return r
}

func (r Rule) GetMessage() string {
	if r.message == "" {
		return r.code
	}
	return r.message
}

func (r Rule) Closure(c RuleClosure) Rule {
	r.closure = c
	return r
}

func (r Rule) GetClosure() RuleClosure {
	return r.closure
}

func (r Rule) Code(c string) Rule {
	r.code = c
	return r
}

func (r Rule) GetCode() string {
	return r.code
}

func (r Rule) Execute(request ValidateRequest) (passes bool, e error) {
	return r.closure(request)
}

// rule with params

func (r RuleWithArgsCtr) Execute(request ValidateRequest) (passes bool, e error) {
	return true, e
}

func (r RuleWithArgs) Args(args ...interface{}) RuleWithArgs {
	r.args = args
	return r
}

func (r RuleWithArgs) GetArgs() []interface{} {
	return r.args
}

func (r RuleWithArgs) Message(m string) RuleWithArgs {
	r.message = m
	return r
}

func (r RuleWithArgs) GetMessage() string {
	return r.message
}

func (r RuleWithArgs) Closure(c RuleClosure) RuleWithArgs {
	r.closure = c
	return r
}

func (r RuleWithArgs) GetClosure() RuleClosure {
	return r.closure
}

func (r RuleWithArgs) Code(c string) RuleWithArgs {
	r.code = c
	return r
}

func (r RuleWithArgs) GetCode() string {
	return r.code
}

// conditional field rules
func (r ConditionalRule) Rules(rules ...interface{}) ConditionalRule {
	r.rules = rules
	return r
}

func (r ConditionalRule) GetRules() []interface{} {
	return r.rules
}

func (r ConditionalRule) FailRules(rules ...interface{}) ConditionalRule {
	r.failRules = rules
	return r
}

func (r ConditionalRule) GetFailRules() []interface{} {
	return r.failRules
}

func (r ConditionalRule) GetApplicableRules(ctx context.Context, value interface{}, params ConditionParams) AcceptableRules {
	if r.condition(ctx, value, params) {
		return r.rules
	}
	return r.failRules
}

func (r ConditionalRule) Condition(cond ConditionalRuleCond) ConditionalRule {
	r.condition = cond
	return r
}

func (r ConditionalRule) GetCondition() ConditionalRuleCond {
	return r.condition
}

// Field rules
func (r FieldRules) StopOnFirstError(v bool) FieldRules {
	r.stopOnFirstError = v
	return r
}

func (r FieldRules) GetStopOnFirstError() bool {
	return r.stopOnFirstError
}

// func (r FieldRules) FieldName(field string) FieldRules {
// 	r.name = field
// 	return r
// }

// func (r FieldRules) GetFieldName() string {
// 	return r.name
// }

func (r FieldRules) Alias(alias string) FieldRules {
	r.alias = alias
	return r
}

func (r FieldRules) GetAlias() string {
	if r.alias != "" {
		return r.alias
	}
	return r.field
}

func (r FieldRules) Field(field string) FieldRules {
	r.field = field
	return r
}

func (r FieldRules) GetField() string {
	return r.field
}

// func (r FieldRules) PathAlias(path string) FieldRules {
// 	r.pathAlias = path
// 	return r
// }

// func (r FieldRules) GetPathAlias() string {
// 	if r.pathAlias != "" {
// 		return r.pathAlias
// 	}
// 	return r.path
// }

func (r FieldRules) Rules(rules AcceptableRules) FieldRules {
	r.rules = rules
	return r
}

func (r FieldRules) GetRules() AcceptableRules {
	return r.rules
}

func (r FieldRules) Messages(messages map[string]interface{}) FieldRules {
	r.messages = utils.ParseMessages(messages)
	return r
}

func (r FieldRules) GetMessages() map[string]string {
	return r.messages
}
