package typetree

import (
	"errors"
	"reflect"
)

// A Value is an arbitrary value
type Value interface {

	// Interface returns the pointer value of the value
	Interface() interface{}

	// Append adds the given value on the array, if it exists and
	// is an array
	Append(value interface{}) error
}

type valueS struct {
	parent interface{}
	v      reflect.Value
	t      reflect.Type
	field  *field
}

func (v *valueS) Interface() interface{} {
	return v.v.Interface()
}

func (v *valueS) Append(value interface{}) error {

	if v.t.Kind() != reflect.Slice {
		return errors.New("value does not support appending")
	}

	v.v = reflect.Append(v.v, reflect.ValueOf(value))

	if v.parent != nil && v.field != nil {
		return v.field.SetValue(v.parent, v.v.Interface())
	}

	return nil
}
