package types

import (
	"context"
	"fmt"
	"strings"

	"github.com/gobigbang/validator/messages"
	"github.com/gobigbang/validator/utils"
)

// Global rule registry (to be used by name in the FieldRules definition)
// ruleCode=>IRule
var Registry = make(RulesRegistry, 0)

// Default validation config
var DefaultValidatorConfig = ValidatorConfig{
	GetValueParams:      utils.DefaultGetValueParams,
	StopOnFirstError:    false,
	StructTag:           "json",
	Translator:          NewTranslator(messages.En),
	ParamsSeparator:     ",",
	NameParamsSeparator: ":",
}

type (
	// Rules to be passed to the validator
	RuleMap struct {
		rules    map[string]FieldRules
		messages map[string]string
	}

	// Configuration object
	ValidatorConfig struct {
		utils.GetValueParams
		StopOnFirstError    bool
		StructTag           string `default:"json"`
		Translator          ITranslator
		ParamsSeparator     string
		NameParamsSeparator string
	}

	// Represents a validation object
	Validator struct {
		rules  RuleMap
		config ValidatorConfig
	}

	// Represents a rule to be executed by the validate method
	preparedRule struct {
		Field   string
		Message string
		Code    string
		Alias   string
		Args    []interface{}
		Closure IValidation
	}
)

// RuleMap methods

// Returns a new RuleMap instance
func NewRuleMap() RuleMap {
	return RuleMap{
		rules: make(map[string]FieldRules),
	}
}

// Set the rules for the map
func (rm RuleMap) Rules(rules map[string]FieldRules) RuleMap {
	rm.rules = rules
	return rm
}

// Obtains the rules of the map
func (rm RuleMap) GetRules() map[string]FieldRules {
	return rm.rules
}

// Set the custom messages for the map (overrides the translator messages and can be overriden by the FieldRule instance)
func (rm RuleMap) Messages(messages map[string]interface{}) RuleMap {
	rm.messages = utils.ParseMessages(messages)
	return rm
}

// Obtains the custom messages
func (rm RuleMap) GetMessage() map[string]string {
	return rm.messages
}

// Executes the prepared rule (implementation of IValidation interface)
func (r preparedRule) Execute(request ValidateRequest) (passes bool, e error) {
	passes, e = r.Closure.Execute(request)
	return
}

// Validator methods

// Creates a new Validator instance
func NewValidator(rules RuleMap) *Validator {
	return &Validator{
		rules:  rules,
		config: DefaultValidatorConfig,
	}
}

// Set the rules to be applied
func (v *Validator) Rules(rules RuleMap) *Validator {
	v.rules = rules
	return v
}

// Obtains the rules
func (v *Validator) GetRules() RuleMap {
	return v.rules
}

// Sets the validator config
func (v *Validator) Config(c ValidatorConfig) *Validator {
	v.config = c
	return v
}

// Obtains the validator config
func (v *Validator) GetConfig() ValidatorConfig {
	return v.config
}

// Gets the value descriptor of the path
func (v *Validator) GetValue(input interface{}, path string) utils.ValueDescriptor {
	return utils.GetValueDescriptorFromPath(input, path, v.config.GetValueParams)
}

// Translates a message using the ITranslator instance
func (v *Validator) Translate(message string, defaultMessage string, args map[string]interface{}) string {
	msg := v.config.Translator.Translate(message, args)
	if msg == "" || msg == message {
		msg = defaultMessage
	}

	return msg
}

// Validates the input against
func (v *Validator) Validate(ctx context.Context, input interface{}) *MessageBag {
	var err *MessageBag
	for _, fieldRules := range v.rules.rules {
		value := v.GetValue(input, fieldRules.GetField())
		applicableRules := v.getApplicableRules(fieldRules, ctx, value)
		for _, r := range applicableRules {
			// switch r.(type) {
			// case Rule:
			// 	rr = r.(Rule)
			// case RuleWithArgs:
			// 	rr = r.(RuleWithArgs)
			// }
			passes, e := r.Execute(ValidateRequest{
				Value:           value.Value,
				Validator:       v,
				Context:         ctx,
				ValueDescriptor: value,
			})

			if !passes {
				if err == nil {
					err = NewMessageBag()
				}

				ve := ValidationError{
					Field:     r.Alias,
					Kind:      value.RKind,
					Rule:      fmt.Sprint(r),
					ErrorCode: r.Code,
					Message:   r.Message,
					Value:     value,
				}

				if e != nil {
					ve.Message = e.Error()
				} else {
					ve.Message = v.Translate(r.Message, ve.Message, map[string]interface{}{
						"Field": r.Alias,
					})
				}

				err.Add(r.Field, ve)

				// if the validation should stop for this field
				if fieldRules.stopOnFirstError {
					break
				}
			}

			// if e != nil {
			// 	if err == nil {
			// 		err = NewMessageBag()
			// 	}

			// 	msg := v.Translate(rr.GetCode(), e.Error(), map[string]interface{}{
			// 		"Field": fieldRules.FieldName,
			// 	})
			// 	ve := ValidationError{
			// 		Field:     fieldRules.GetAlias(),
			// 		Kind:      value.RKind,
			// 		Rule:      fmt.Sprint(r),
			// 		ErrorCode: rr.GetCode(),
			// 		Message:   msg,
			// 		Value:     value,
			// 	}
			// 	err.Add(fieldRules.GetPath(), ve)
			// 	// Stop if error
			// 	if fieldRules.stopOnFirstError {
			// 		break
			// 	}
			// }
		}
	}
	return err
}

// Obtains a rule from the global registry
func (v *Validator) getRule(r interface{}) IValidation {
	var ret IValidation
	switch r.(type) {
	case string:
		s, ok := Registry[r.(string)]
		if ok {
			return s
		}
	}
	return ret
}

// Parses the rule name (splits the params and formats the name)
func (v *Validator) parseRuleName(r string) (rule string, params []string) {
	parts := strings.Split(r, v.config.NameParamsSeparator)
	rule = parts[0]
	if len(parts) > 1 {
		params = strings.Split(parts[1], v.config.ParamsSeparator)
	}
	return
}

// Obtains the applicable rules for a field
// All the conditional rules conditions will be tested to determine if should be added or not
func (v *Validator) getApplicableRules(rules FieldRules, ctx context.Context, value interface{}) []preparedRule {
	applicableRules := make([]preparedRule, 0)
	for _, rule := range rules.GetRules() {
		pr := preparedRule{
			Field: rules.GetField(),
			Alias: rules.GetAlias(),
		}

		// when the field is a simple var, the field will be empty
		if pr.Field == "" {
			pr.Field = pr.Alias
		}

		switch rule.(type) {
		// case ConditionalRule:
		// 	params := ConditionParams{
		// 		Validator: v,
		// 	}
		// 	rArr := rule.(ConditionalRule).GetApplicableRules(ctx, value, params)
		// 	for _, candidate := range rArr {
		// 		r := v.getRule(candidate)
		// 		if r != nil {
		// 			pr.Closure = r
		// 		}
		// 	}
		// case Rule:
		// 	applicableRules = append(applicableRules, rule.(Rule))
		// case RuleWithArgs:
		// 	applicableRules = append(applicableRules, rule.(RuleWithArgs))
		case string:
			rName, _ := v.parseRuleName(rule.(string))
			r := v.getRule(rName)

			// Get the message from the validation rules
			pr.Message = v.rules.messages[rName]

			if r != nil {
				rr := r.(IRule)

				switch r.(type) {
				// case RuleWithArgsCtr:
				// 	var iargs []interface{}
				// 	for _, a := range args {
				// 		iargs = append(iargs, a)
				// 	}
				// 	applicableRules = append(applicableRules, r.(RuleWithArgsCtr)(iargs...))
				default:
					pr.Code = rr.GetCode()
					if pr.Message == "" {
						pr.Message = rr.GetMessage()
					}
					pr.Closure = rr.GetClosure()
					if msg := rules.messages[pr.Code]; msg != "" {
						pr.Message = msg
					}
					applicableRules = append(applicableRules, pr)
				}

			}
		}
	}
	return applicableRules
}
