package utils

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type GetValueParams struct {
	StructTag              string `default:"json"`
	PathSeparator          string `default:"."`
	ArrayWildcard          string `default:"*"`
	ArrayIndexesRegexp     regexp.Regexp
	ArraySingleIndexRegexp regexp.Regexp
}

type (
	// Container for the reflection values of a field or value
	ValueDescriptor struct {
		//the kind of the value
		RKind reflect.Kind
		// the type of the value
		RType reflect.Type
		// the reflect.Value
		RValue reflect.Value
		// the value
		Value interface{}
		// is simple value? (not an array, slice, struct or map)
		IsSimple bool
		// the value is a ponter?
		IsPtr bool
	}
)

var DefaultGetValueParams = GetValueParams{
	StructTag:              "json",
	PathSeparator:          ".",
	ArrayWildcard:          "*",
	ArrayIndexesRegexp:     *regexp.MustCompile(`\[([0-9\,\*]+)\]`),
	ArraySingleIndexRegexp: *regexp.MustCompile(`\{([0-9]+)\}`),
}

/*
Returns a value descriptor for a value based on path
*/
func GetValueDescriptorFromPath(value interface{}, path string, config GetValueParams) ValueDescriptor {
	vdescriptor := GetValueDescriptor(value)

	if path == "" {
		return vdescriptor
	}

	segments := strings.Split(path, config.PathSeparator)
	if len(segments) == 0 || vdescriptor.RKind == reflect.Invalid {
		return vdescriptor
	}

	field := segments[0]
	switch vdescriptor.RKind {
	case reflect.Struct:
		var fvalue reflect.Value
		var found bool = false
		if config.StructTag != "" {
			fields := vdescriptor.RType.NumField()

			for i := 0; i < fields; i++ {
				sField := vdescriptor.RType.Field(i)
				tag := ""
				if config.StructTag != "" {
					tag = strings.Split(sField.Tag.Get(config.StructTag), ",")[0]
				}
				if tag == "" {
					tag = sField.Name
				}
				if tag == field {
					found = true
					fvalue = vdescriptor.RValue.Field(i)
				}
			}
		}
		if !found {
			fvalue = vdescriptor.RValue.FieldByName(field)
		}
		if !fvalue.IsValid() {
			return CreateInvalidDescriptor()
		}
		v := fvalue.Interface()
		if len(segments) > 1 {
			return GetValueDescriptorFromPath(v, strings.Join(segments[1:], config.PathSeparator), config)
		} else {
			return GetValueDescriptor(v)
		}
	case reflect.Array, reflect.Slice:
		if len(segments) > 0 {
			istr := segments[0]
			index := config.ArraySingleIndexRegexp.FindStringSubmatch(istr)
			p := strings.Join(segments[1:], config.PathSeparator)
			if len(index) == 2 {
				idx := index[1]
				i, err := strconv.Atoi(idx)
				if err == nil && i < vdescriptor.RValue.Len() {
					vitem := vdescriptor.RValue.Index(i)
					if vitem.IsValid() {
						vdescriptor = GetValueDescriptorFromPath(vitem.Interface(), p, config)
					}
				}
			} else {
				indexes := config.ArrayIndexesRegexp.FindStringSubmatch(istr)
				if len(indexes) == 2 {
					p := strings.Join(segments[1:], config.PathSeparator)
					items := make([]interface{}, 0)
					if strings.Contains(indexes[1], config.ArrayWildcard) {
						for i := 0; i < vdescriptor.RValue.Len(); i++ {
							vitem := vdescriptor.RValue.Index(i)
							if vitem.IsValid() {
								currItems := GetValueDescriptorFromPath(vitem.Interface(), p, config).Value
								items = append(items, currItems)
								// d := GetValueDescriptor(currItems)
								// if d.RKind == reflect.Array || d.RKind == reflect.Slice {
								// 	for i := 0; i < d.RValue.Len(); i++ {
								// 		vitem := d.RValue.Index(i)
								// 		if vitem.IsValid() {
								// 			items = append(items, vitem.Interface())
								// 		}
								// 	}
								// } else {
								// 	items = append(items, currItems)
								// }
							}

						}
					} else {
						searchIndexes := strings.Split(indexes[1], ",")
						for _, si := range searchIndexes {
							i, err := strconv.Atoi(si)
							if err == nil && i < vdescriptor.RValue.Len() {
								vitem := vdescriptor.RValue.Index(i)
								if vitem.IsValid() {
									currItems := GetValueDescriptorFromPath(vitem.Interface(), p, config).Value
									items = append(items, currItems)
									// d := GetValueDescriptor(currItems)
									// if d.RKind == reflect.Array || d.RKind == reflect.Slice {
									// 	for i := 0; i < d.RValue.Len(); i++ {
									// 		vitem := d.RValue.Index(i)
									// 		if vitem.IsValid() {
									// 			items = append(items, vitem.Interface())
									// 		}
									// 	}
									// } else {
									// 	items = append(items, currItems)
									// }
								}
							}
						}
					}
					vdescriptor.Value = items
				}
			}

		}
		return vdescriptor
	case reflect.Map:
		mapInput, ok := value.(map[string]interface{})
		if !ok {
			return CreateInvalidDescriptor()
		}
		v, ok := mapInput[field]
		if !ok {
			return CreateInvalidDescriptor()
		}
		if len(segments) > 1 {
			return GetValueDescriptorFromPath(v, strings.Join(segments[1:], config.PathSeparator), config)
		} else {
			return GetValueDescriptor(v)
		}
	default:
		return vdescriptor
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

func GetValueDescriptor(value interface{}) ValueDescriptor {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	p := false
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
		p = true
	}

	var i interface{} = value
	// if (v.IsValid()) {
	// 	i = v.Interface()
	// }

	return ValueDescriptor{
		RValue:   v,
		RType:    t,
		RKind:    k,
		Value:    i,
		IsPtr:    p,
		IsSimple: IsSimpleValue(k),
	}
}
