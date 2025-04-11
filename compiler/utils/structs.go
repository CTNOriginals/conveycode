package utils

import (
	"reflect"
)

func StructKeys(obj any) (keys []string) {
	var val = reflect.Indirect(reflect.ValueOf(obj))
	var length = val.Type().NumField()
	keys = make([]string, length)

	for i := range length {
		keys[i] = val.Type().Field(i).Name
	}

	return keys
}
func StructValues(obj any) (vals []any) {
	var val = reflect.Indirect(reflect.ValueOf(obj))
	var length = val.Type().NumField()
	vals = make([]any, length)

	for i := range length {
		vals[i] = val.Field(i)
	}

	return vals
}
