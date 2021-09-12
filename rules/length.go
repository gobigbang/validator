package rules

import (
	"context"
	"reflect"
	"strconv"

	"github.com/gobigbang/validator/messages"
	"github.com/gobigbang/validator/types"
)

func MinLength(ctx context.Context, params types.RuleCheckerParameters) error {
	var err error
	if params.Value != nil && len(params.Args) == 1 {
		var passes = false
		min, perr := strconv.Atoi(params.Args[0])
		if perr == nil {
			d := types.GetDescriptor(params.Value)
			switch d.RKind {
			case reflect.Slice, reflect.Array, reflect.String, reflect.Map:
				passes = d.RValue.Len() >= min
			}

			if !passes {
				t := "string"
				switch d.RKind {
				case reflect.Array, reflect.Slice:
					t = "array"
				case reflect.Map:
					t = "map"
				}
				err = types.NewValidationError("min_length", messages.GetTypedMessage(params.Messages, "min_length", t), params)
			}
		}
	}
	return err
}

func MaxLength(ctx context.Context, params types.RuleCheckerParameters) error {
	var err error
	if params.Value != nil {
		var passes = false
		max, perr := strconv.Atoi(params.Args[0])
		if perr == nil {
			d := types.GetDescriptor(params.Value)
			switch d.RKind {
			case reflect.Slice, reflect.Array, reflect.String, reflect.Map:
				passes = d.RValue.Len() <= max
			default:
				passes = true
			}

			if !passes {
				t := "string"
				switch d.RKind {
				case reflect.Array, reflect.Slice:
					t = "array"
				case reflect.Map:
					t = "map"
				}
				err = types.NewValidationError("max_length", messages.GetTypedMessage(params.Messages, "max_length", t), params)
			}
		}
	}
	return err
}
