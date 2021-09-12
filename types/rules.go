package types

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gobigbang/validator/messages"
)

type (
	Rule struct {
		Name    string
		Message string
	}

	// Rules
	RulesArr []interface{}

	// Field rule definition
	FieldRules struct {
		FieldName string
		Rules     RulesArr
	}

	// Field rules fieldName:FieldRuleDefinition
	FieldRulesType map[string]FieldRules

	ArrayRules struct {
		Rules map[string]ArrayRulesItem
	}

	ArrayRulesItem struct {
		Key   string
		Rules []interface{}
	}

	// Rule map
	MapRules struct {
		FieldRules       map[string]FieldRules
		ExtraKeysAllowed []string
	}

	// Parameters to be passed to the rule function from the validator
	RuleCheckerParameters struct {
		Field        string
		Value        interface{}
		Values       interface{}
		Args         []string
		ErrorMessage string
		Descriptor   ValueDescriptor
		Messages     map[string]interface{}
	}
	// Rule check function definition
	RuleCheckerFunc func(ctx context.Context, params RuleCheckerParameters) error

	ParsedRule struct {
		Checker RuleCheckerFunc
		Args    []string
	}

	ConditionalRuleFunc func(ctx context.Context, value interface{}, values interface{}) bool
	ConditionalRule     struct {
		Cond  ConditionalRuleFunc
		Rules RulesArr
	}

	Validatable interface {
		Validate(ctx context.Context, input interface{}, config ValidatorConfig) error
	}
)

func (fr FieldRules) Validate(ctx context.Context, input interface{}, config ValidatorConfig) error {
	var err error
	fieldValue := GetValueFromPath(fr.FieldName, input, config)
	for _, r := range fr.Rules {
		var ruleFunc ParsedRule
		var vResult error
		var errMessage string
		switch r.(type) {
		case string:
			ruleFunc = getRuleFunc(r)
			if ruleFunc.Checker == nil {
				vResult = NewInternalError("invalid_rule", messages.GetMessage(config.Messages, "invalid_rule"), map[string]interface{}{
					"Rule": fmt.Sprint(r),
				})
			}
		case Rule:
			rr := r.(Rule)
			ruleFunc = getRuleFunc(rr.Name)
			errMessage = rr.Message
		case MapRules:
			rm := r.(MapRules)
			vResult = rm.Validate(ctx, fieldValue, config)
		case ArrayRules:
			ar := r.(ArrayRules)
			vResult = ar.Validate(ctx, fieldValue, config)
		case ConditionalRule:
			cr := r.(ConditionalRule)
			if cr.Cond(ctx, fieldValue, input) {
				f := FieldRules{
					Rules:     cr.Rules,
					FieldName: fr.FieldName,
				}
				vResult = Validate(ctx, fieldValue, f, config)
			}
		default:
			vResult = NewInternalError("invalid_rule", messages.GetMessage(config.Messages, "invalid_rule"), map[string]interface{}{
				"Rule": fmt.Sprint(r),
			})
		}

		if ruleFunc.Checker != nil {
			d := GetDescriptor(fieldValue)
			params := RuleCheckerParameters{
				Field:        fr.FieldName,
				Value:        d.Value,
				Descriptor:   d,
				Args:         ruleFunc.Args,
				ErrorMessage: errMessage,
				Messages:     config.Messages,
			}
			vResult = ruleFunc.Checker(ctx, params)

		}

		if vResult != nil {
			if err == nil {
				err = make(ValidationErrors, 0)
			}
			verr := err.(ValidationErrors)
			if e, ok := vResult.(ValidationErrors); ok {
				err = append(verr, e...)
			} else {
				err = append(verr, vResult)
			}
		}
	}
	return err
}

func (rm MapRules) Validate(ctx context.Context, input interface{}, config ValidatorConfig) error {
	var err error
	for fieldName, fieldRules := range rm.FieldRules {
		vErr := fieldRules.Validate(ctx, input, config)
		if vErr != nil {
			if err == nil {
				err = make(ValidationErrorMap)
			}
			em := err.(ValidationErrorMap)
			if _, ok := em[fieldName]; !ok {
				em[fieldName] = make(ValidationErrors, 0)
			}
			switch vErr.(type) {
			case ValidationErrors:
				em[fieldName] = append(em[fieldName], vErr.(ValidationErrors)...)
			default:
				em[fieldName] = append(em[fieldName], vErr)
			}

		}
	}
	return err
}

func (ar ArrayRules) Validate(ctx context.Context, input interface{}, config ValidatorConfig) error {
	var err error
	d := GetDescriptor(input)
	if d.RKind == reflect.Array || d.RKind == reflect.Slice {
		var vResult error
		var applyToAll []interface{}
		for arrKey, arrRules := range ar.Rules {
			// all values
			if arrKey == "*" {
				applyToAll = arrRules.Rules
			} else {
				var subErrors ValidationErrorMap
				i, converr := strconv.Atoi(arrKey)
				if converr == nil {
					var sv interface{} = nil
					if d.RValue.Len() >= i+1 {
						iv := d.RValue.Index(i)
						if (iv != reflect.Value{}) {
							sv = iv.Interface()
						}
					}
					applyRules := arrRules.Rules
					applyRules = append(applyRules, applyToAll...)
					for _, r := range applyRules {
						e := Validate(ctx, sv, r, config)
						if e != nil {
							if subErrors == nil {
								subErrors = make(ValidationErrorMap)
							}
							if emap, ok := e.(ValidationErrors); ok {
								if subErrors[fmt.Sprint(i)] == nil {
									subErrors[fmt.Sprint(i)] = make(ValidationErrors, 0)
								}
								for _, emValue := range emap {
									subErrors[fmt.Sprint(i)] = append(subErrors[fmt.Sprint(i)], emValue)
								}
							}
						}
					}
				}
				if subErrors != nil {
					if vResult == nil {
						vResult = make(ValidationErrorMap)
					}

					for k, v := range subErrors {
						vResult.(ValidationErrorMap)[k] = v
					}

				}

			}
			err = vResult
		}

	}
	return err
}

func (cr ConditionalRule) Validate(ctx context.Context, input interface{}, config ValidatorConfig) error {
	var err error
	if cr.Cond(ctx, input, input) {
		err = Validate(ctx, input, cr.Rules, config)
	}
	return err
}

func (ra RulesArr) Validate(ctx context.Context, input interface{}, config ValidatorConfig) error {
	var err error
	for _, r := range ra {
		err = Validate(ctx, input, r, config)
	}
	return err
}
