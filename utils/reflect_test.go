package utils_test

import (
	"fmt"
	"testing"

	"github.com/gobigbang/validator/utils"
)

type S1 struct {
	A string `json:"a"`
	B string `json:"b"`
	C int    `json:"c"`
}

func TestReflect(t *testing.T) {
	t.Run("Simple value", func(st *testing.T) {
		input := "Hello world"
		d := utils.GetValueDescriptorFromPath(input, "hello", utils.DefaultGetValueParams)
		if d.Value != input {
			st.Error(fmt.Sprint("Value is: ", d.Value))
		}
	})

	t.Run("Map value", func(st *testing.T) {
		input := map[string]interface{}{
			"a": "Item a",
			"b": "Item b",
			"c": map[string]interface{}{
				"a": "Item c_a",
			},
			"d": &S1{
				A: "hello world a",
				B: "hello world b",
			},
		}
		d := utils.GetValueDescriptorFromPath(input, "c.a", utils.DefaultGetValueParams)
		if d.Value != "Item c_a" {
			st.Error(fmt.Sprint("Value is: ", d.Value))
		}

		d1 := utils.GetValueDescriptorFromPath(input, "d", utils.DefaultGetValueParams)
		d2 := utils.GetValueDescriptorFromPath(input, "d", utils.DefaultGetValueParams)
		i := d1.Value.(*S1)
		i.A = "Pepe"
		if d1.Value.(*S1).A != d2.Value.(*S1).A {
			st.Error("Values not maches", d1.Value, d2.Value)
		} else {
			st.Log(d1.Value.(*S1).A, i.A)
		}
	})

	t.Run("Array value", func(st *testing.T) {
		var input []map[string]interface{}

		for i := 0; i < 10; i++ {
			input = append(input, map[string]interface{}{
				"a": "Item a" + fmt.Sprint(i),
				"b": "Item b" + fmt.Sprint(i),
				"c": map[string]interface{}{
					"a": []map[string]interface{}{
						{"a": "Item c_a_a_0_" + fmt.Sprint(i)},
						{"a": "Item c_a_a_1_" + fmt.Sprint(i)},
					},
				},
			})
		}

		d := utils.GetValueDescriptorFromPath(input, "[*].c.a.{0}", utils.DefaultGetValueParams)
		if d.RValue.Len() != 10 {
			t.Error("Invalid len, expeted 10")
		} else {
			t.Log(utils.JsonPrettyPrint(d))
		}
	})
}
