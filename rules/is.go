package rules

// import (
// 	"context"
// 	"errors"
// 	"reflect"

// 	"github.com/gobigbang/validator/types"
// )

// var IsString = types.Closure(func(ctx context.Context, value interface{}, params types.ValidateRequest) (e error) {
// 	passes := value == nil || params.ValueDescriptor.RKind == reflect.String
// 	if !passes {
// 		e = errors.New(params.Message)
// 	}
// 	return
// }).Message("The field must be a string.").Code("is_string")

// var IsNumeric = types.Closure(func(ctx context.Context, value interface{}, params types.ValidateRequest) (e error) {
// 	passes := value == nil

// 	if params.ValueDescriptor.RKind >= reflect.Int && params.ValueDescriptor.RKind <= reflect.Complex128 {
// 		passes = true
// 	}

// 	if !passes {
// 		e = errors.New(params.Message)
// 	}
// 	return
// }).Message("The field must be a number.").Code("is_numeric")
