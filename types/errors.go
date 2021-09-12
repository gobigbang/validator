package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gobigbang/validator/messages"
)

type (
	// Represents a validation error on field
	ValidationError struct {
		Field     string
		Kind      reflect.Kind
		Rule      string
		Value     interface{}
		ErrorCode string
		Message   string
	}

	InternalError struct {
		ErrorCode string
		Message   string
		Params    interface{}
	}

	// Map errors to field
	ValidationErrorMap map[string]ValidationErrors
	ValidationErrors   []error
)

func (e InternalError) Error() string {
	return e.Message
}

func (e ValidationErrors) Error() string {
	count := len(e)
	if count > 0 {
		allErrors := make([]string, count)
		if count > 0 {
			for i, v := range e {
				allErrors[i] = v.Error()
			}
		}
		return fmt.Sprint(allErrors)
	}
	return ""
}

func (ve ValidationError) Error() string {
	return ve.Message
}

func (em ValidationErrorMap) Error() string {

	count := len(em)
	if count > 0 {
		allErrors := make([]string, 0)
		for k, v := range em {
			errList := make([]string, 0)
			for _, e := range v {
				if vErr, ok := e.(ValidationError); ok {
					errList = append(errList, vErr.Error())
				}
			}
			allErrors = append(allErrors, fmt.Sprintf("%s:%s", k, errList))
		}

		return fmt.Sprint(allErrors)
	}
	return ""
}

// Creates a new validation error
func NewValidationError(rule string, defaultMessage string, params RuleCheckerParameters) ValidationError {
	message := defaultMessage
	if params.ErrorMessage != "" {
		message = params.ErrorMessage
	}
	message = ParseMessage(message, params)
	vErr := ValidationError{
		Field:     params.Field,
		Kind:      params.Descriptor.RKind,
		Value:     params.Value,
		ErrorCode: fmt.Sprint("validation", "_", rule),
		Rule:      rule,
		Message:   message,
	}
	return vErr
}

// Creates a new internal error
func NewInternalError(code string, defaultMessage string, params map[string]interface{}) InternalError {
	message := defaultMessage
	if message == "" {
		message = messages.Messages[code].(string)
		if message == "" {
			message = "error_" + code
		}
	}
	if message != "" {
		args := make([]string, 0)
		if params != nil {
			for k, v := range params {
				args = append(args, "{"+k+"}")
				args = append(args, fmt.Sprint(v))
			}
		}
		message = strings.NewReplacer(args...).Replace(message)
	}
	return InternalError{
		ErrorCode: code,
		Message:   message,
		Params:    params,
	}
}

func ParseMessage(format string, params RuleCheckerParameters) string {
	args := make([]string, 0)

	args = append(args, "{Field}", params.Field, "{Value}", fmt.Sprint(params.Value))
	i := 0
	for _, v := range params.Args {
		args = append(args, "{Arg"+fmt.Sprint(i)+"}", v)
	}
	return strings.NewReplacer(args...).Replace(format)
}
