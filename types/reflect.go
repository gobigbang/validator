package types

import (
	"reflect"
	"strings"

	"github.com/gobigbang/validator/utils"
)

type (
	// Container for the reflection values of a field or value
	ValueDescriptor struct {
		RKind    reflect.Kind
		RType    reflect.Type
		RValue   reflect.Value
		Value    interface{}
		IsSimple bool
	}

	// Value descriptor list
	ValueDescriptors []ValueDescriptor
)

// Obtains the reflect values of the input
func GetDescriptor(input interface{}) ValueDescriptor {
	t := reflect.TypeOf(input)
	v := reflect.ValueOf(input)
	var k reflect.Kind
	if t != nil {
		k = t.Kind()
	} else {
		k = reflect.Invalid
	}

	// the input is a ponter, we obtain the underlying type
	if k == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
		k = t.Kind()
	}

	var i interface{} = input
	if (v != reflect.Value{}) {
		i = v.Interface()
	}

	return ValueDescriptor{
		RValue:   v,
		RType:    t,
		RKind:    k,
		Value:    i,
		IsSimple: IsSimpleValue(k),
	}
}

// Determines if the field is a basic type or a map, array, slice or struct
func IsSimpleValue(kind reflect.Kind) bool {
	switch kind {
	case reflect.Struct, reflect.Array, reflect.Map, reflect.Slice:
		return false
	default:
		return true
	}
}

func CreateInvalidDescriptor() ValueDescriptor {
	return ValueDescriptor{
		RKind:  reflect.Invalid,
		RValue: reflect.Value{},
		RType:  nil,
		Value:  nil,
	}
}

func GetValueFromPath(path string, input interface{}, config ValidatorConfig) interface{} {
	if path == "" {
		return input
	}
	d := GetDescriptor(input)
	segments := strings.Split(path, ".")
	if len(segments) == 0 || d.RKind == reflect.Invalid {
		return nil
	}

	fieldName := segments[0]

	switch d.RKind {
	case reflect.Struct:
		var f reflect.Value
		var found bool = false
		if config.StructTag != "" {
			fields := d.RType.NumField()

			for i := 0; i < fields; i++ {
				field := d.RType.Field(i)
				tag := ""
				if config.StructTag != "" {
					tag = strings.Split(field.Tag.Get(config.StructTag), ",")[0]
				}
				if tag == "" {
					tag = field.Name
				}
				if tag == fieldName {
					found = true
					f = d.RValue.Field(i)
				}
			}
		}
		if !found {
			f = d.RValue.FieldByName(fieldName)
		}
		if !f.IsValid() {
			return nil
		}
		v := f.Interface()
		if len(segments) > 1 {
			return GetValueFromPath(strings.Join(segments[1:], "."), v, config)
		} else {
			return v
		}
	case reflect.Array, reflect.Slice:
		//arrInput := ConvertToInterfaceSlice(input, d)
		//return arrInput
		return input
	case reflect.Map:
		mapInput, ok := input.(map[string]interface{})
		if !ok {
			return nil
		}
		v, ok := mapInput[fieldName]
		if !ok {
			return nil
		}
		if len(segments) > 1 {
			return GetValueFromPath(strings.Join(segments[1:], "."), v, config)
		} else {
			return v
		}
	default:
		return input
	}
}

func GetKeys(input interface{}, structTag string) []string {
	ret := make([]string, 0)
	d := GetDescriptor(input)
	switch d.RKind {
	case reflect.Map:
		keys := d.RValue.MapKeys()
		for _, k := range keys {
			ret = append(ret, k.String())
		}
	case reflect.Struct:
		fields := d.RType.NumField()

		for i := 0; i < fields; i++ {
			field := d.RType.Field(i)
			tag := ""
			if structTag != "" {
				tag = strings.Split(field.Tag.Get(structTag), ",")[0]
			}
			if tag == "" {
				tag = field.Name
			}
			ret = append(ret, tag)
		}
	}

	return ret
}

func MakeCompleteKey(keys ...string) string {
	return strings.Join(utils.FilterStrings(keys, func(s string) bool {
		return s != ""
	}), ".")
}
