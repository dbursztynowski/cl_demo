package tools

import (
	"reflect"
)

type MyReflectValue struct {
	reflect.Value
}

type Invalid struct {
	int
}

func (v MyReflectValue) XFieldByName(name string) MyReflectValue {
	if v.Kind() == reflect.Struct {
		return MyReflectValue{v.FieldByName(name)}
	} else {
		invalid := Invalid{0}
		reflectInvalid := reflect.ValueOf(invalid)
		return MyReflectValue{reflectInvalid}
	}
}

func (v MyReflectValue) XSetString(value string) {
	if v.Kind() == reflect.String {
		v.SetString(value)
	}
}

func (v MyReflectValue) XSet(value reflect.Value) {
	if v.Kind() != 0 {
		v.Set(value)
	}
}
