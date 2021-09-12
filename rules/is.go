package rules

import (
	"context"
	"reflect"

	"github.com/gobigbang/validator/messages"
	"github.com/gobigbang/validator/types"
)

func IsArray(ctx context.Context, params types.RuleCheckerParameters) error {
	var err error
	passes := params.Value == nil || params.Descriptor.RKind == reflect.Array || params.Descriptor.RKind == reflect.Slice
	if !passes {
		err = types.NewValidationError("is_array", messages.GetMessage(params.Messages, "is_array"), params)
	}
	return err
}

func IsString(ctx context.Context, params types.RuleCheckerParameters) error {
	var err error

	passes := params.Value == nil || params.Descriptor.RKind == reflect.String
	if !passes {
		err = types.NewValidationError("is_string", messages.GetMessage(params.Messages, "is_string"), params)
	}
	return err
}

func IsInteger(ctx context.Context, params types.RuleCheckerParameters) error {
	var err error

	passes := params.Value == nil || params.Descriptor.RKind == reflect.Int || params.Descriptor.RKind == reflect.Int16 || params.Descriptor.RKind == reflect.Int8 || params.Descriptor.RKind == reflect.Int32 || params.Descriptor.RKind == reflect.Int64
	if !passes {
		err = types.NewValidationError("is_integer", messages.GetMessage(params.Messages, "is_integer"), params)
	}
	return err
}

func IsMap(ctx context.Context, params types.RuleCheckerParameters) error {
	var err error

	passes := params.Value == nil || params.Descriptor.RKind == reflect.Map
	if !passes {
		err = types.NewValidationError("is_map", messages.GetMessage(params.Messages, "is_map"), params)
	}
	return err
}

func IsStruct(ctx context.Context, params types.RuleCheckerParameters) error {
	var err error

	passes := params.Value == nil || params.Descriptor.RKind == reflect.Struct
	if !passes {
		err = types.NewValidationError("is_struct", messages.GetMessage(params.Messages, "is_struct"), params)
	}
	return err
}
