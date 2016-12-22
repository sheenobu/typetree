package typetree

import (
	"reflect"
	"testing"
)

func TestValueSimple(t *testing.T) {

	var str = "hello"
	var v valueS
	v.v = reflect.ValueOf(str)
	v.t = reflect.TypeOf(str)

	if str != v.Interface().(string) {
		t.Error("Failed the simplest test")
	}

}

func TestValueAppend(t *testing.T) {

	var ix []string
	var v valueS
	v.v = reflect.ValueOf(ix)
	v.t = reflect.TypeOf(ix)

	err := v.Append("hello")
	if err != nil {
		t.Errorf("Err: %v", err)
	}

	if len(ix) != 0 {
		t.Errorf("original list should be len 0")
	}

	ix2 := v.Interface().([]string)

	if len(ix2) != 1 {
		t.Errorf("new list should be len 1")
	}
}

func TestValueAppendError(t *testing.T) {

	var ix string
	var v valueS
	v.v = reflect.ValueOf(ix)
	v.t = reflect.TypeOf(ix)

	err := v.Append("hello")
	if err == nil {
		t.Errorf("Expected error")
	}

}
