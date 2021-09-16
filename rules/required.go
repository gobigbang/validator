package rules

import "github.com/gobigbang/validator/types"

var Required = types.NewRule(func(request types.ValidateRequest) (passes bool, e error) {
	passes = request.Value != nil && request.ValueDescriptor.RValue.IsValid()
	return
}).Code("required")

var NotZero = types.NewRule(func(request types.ValidateRequest) (passes bool, e error) {
	passes = request.Value != nil && request.ValueDescriptor.RValue.IsValid() && !request.ValueDescriptor.RValue.IsZero()
	return
}).Code("not_zero")
