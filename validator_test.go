package validator_test

import (
	"context"
	"testing"

	"github.com/gobigbang/validator"
)

func TestValidator(t *testing.T) {

	t.Run("Simple var validation", func(t *testing.T) {
		input := ""
		rules := validator.Var("Input Var", "required", "not_zero").Messages(map[string]interface{}{
			"required": "El campo {Field} no puede ser nulo",
			"not_zero": "El campo {Field} no puede ser 0",
		})
		v := validator.New(rules)
		e := v.Validate(context.Background(), input)
		if e != nil {
			t.Error(e)
		}
	})

	// t.Run("Map validator", func(t *testing.T) {
	// 	rules:=validator.Rules()
	// })

	// rules := validator.Rules(
	// 	validator.Field("a", `required`, `not_zero`).Messages(map[string]string{
	// 		"required": `El campo "{Field}" tiene que estar`,
	// 		"not_zero": `El campo "{Field}" no puede ser 0 value`,
	// 	}).Alias("field01"),
	// // validator.Field("a", rules.Required.Message("The field coso"), rules.In(1, 2, 3), rules.NotIn(1)),
	// // validator.Field("b", "required", validator.Conditional(func(ctx context.Context, value interface{}, params types.ConditionParams) bool {
	// // 	return false
	// // }).Rules("string").FailRules("numeric")),
	// // validator.Field("c", "required", "string", validator.Closure(func(ctx context.Context, value interface{}, params types.ValidateRequest) (e error) {
	// // 	var err error
	// // 	if value == nil {
	// // 		return errors.New("Coso")
	// // 	}
	// // 	return err
	// // })),
	// // validator.Field("d", "not_zero"),
	// )

	// input := map[string]interface{}{
	// 	"a": 0,
	// 	"b": false,
	// 	"c": 0,
	// 	//"d": true,
	// }

	// v := validator.New(rules)
	// err := v.Validate(context.Background(), input)

	// t.Error(utils.JsonPrettyPrint(err.(*types.MessageBag).Messages()))
}
