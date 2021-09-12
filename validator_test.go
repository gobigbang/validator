package validator_test

import (
	"context"
	"testing"
	"time"

	"github.com/gobigbang/validator"
	"github.com/gobigbang/validator/utils"
)

type Struct01 struct {
	A string `json:"a"`
	B int
	C time.Time
}

func TestSimpleValidation(t *testing.T) {
	rules := validator.FieldRules("var01", "string", "required", "min_length:8")
	v := validator.Default(rules)
	var input = "Hello World"
	err := v.Validate(context.Background(), input)
	if err != nil {
		t.Error(err)
	}
}

func TestMapValidation(t *testing.T) {
	// rules := validator.MapRules(
	// 	// validator.FieldRules("a", "string", "required", "max_length:4", "min_length:90"),
	// 	// validator.FieldRules("b", "integer", "required"),
	// 	// validator.FieldRules("c", "array", "min_length:3", validator.ArrayRules(
	// 	// 	validator.ArrayItemRules("*", validator.FieldRules("c_c", "required")),
	// 	// 	validator.ArrayItemRules("0", validator.FieldRules("a", "required")),
	// 	// )),
	// 	validator.FieldRules("d", "struct", validator.MapRules(
	// 		validator.FieldRules("a", "required", "max_length:90", validator.Conditional(func(ctx context.Context, value, values interface{}) bool {
	// 			return true
	// 		}, "min_length:80")),
	// 	)),
	// )
	rules := validator.Conditional(func(ctx context.Context, value, values interface{}) bool {
		return true
	}, validator.MapRules(
		validator.FieldRules("a", "string", "required", "max_length:4", "min_length:90"),
		validator.FieldRules("b", "integer", "required", validator.Conditional(func(ctx context.Context, value, values interface{}) bool {
			return true
		}), "min:1"),
	))
	v := validator.Default(rules)
	var input = map[string]interface{}{
		"a": "Hello world",
		"b": 9.0,
		"c": []map[string]interface{}{
			{
				"c_a": "Hello c_a",
				"c_b": 2,
			},
		},
		"d": &Struct01{
			A: "structfielda",
		},
	}
	err := v.Validate(context.Background(), input)
	if err != nil {
		t.Error(utils.JsonPrettyPrint(err))
	}
}
