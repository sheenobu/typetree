package typetree

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// this is the actual type tree implementation

type tree struct {
	Type   reflect.Type
	fields map[string]*field
}

type field struct {
	Identifier identifier
	Type       reflect.Type
	Index      int

	tree *tree // a field itself could have a subtree
}

func (f *field) SetValue(dest interface{}, val interface{}) error {

	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}

	rf := v.Field(f.Index)

	if !rf.CanSet() {
		return errors.Errorf("Field is not settable: %T %v", dest, rf.Type())
	}

	rv := reflect.ValueOf(val)

	assignable, ptrLeft, ptrRight := assignableToCheck(f.Type, rv)
	if !assignable {
		return errors.Errorf("Value '%v' is not assignable to field", val)
	}

	if ptrLeft {
		rf.Set(reflect.New(rv.Type()))
		rf.Elem().Set(rv)
	} else if ptrRight {
		rf.Set(rv.Elem())
	} else {
		rf.Set(rv)
	}

	return nil
}

func (f *field) GetValue(dest interface{}) (interface{}, error) {
	rf := reflect.ValueOf(dest)

	if rf.Kind() == reflect.Ptr {
		rf = rf.Elem()
	}

	rf = rf.Field(f.Index)

	// rf.Interface() will return a non-nil value ALWAYS, even if IsNil is true.
	if rf.Kind() == reflect.Ptr && rf.IsNil() {
		return nil, nil
	}

	i := rf.Interface()

	return i, nil
}

func assignableToCheck(left reflect.Type, rightValue reflect.Value) (assignable bool, ptrLeft bool, ptrRight bool) {

	assignable = left.AssignableTo(rightValue.Type())
	if assignable {
		return
	}

	assignable = left.AssignableTo(reflect.PtrTo(rightValue.Type()))
	if assignable {
		ptrLeft = true
		return
	}

	assignable = rightValue.Type().AssignableTo(reflect.PtrTo(left))
	if assignable {
		ptrRight = true
		return
	}

	return
}

type identifier struct {
	Name string
	Tag  string
}

func buildTree(tag string, t reflect.Type) (*tree, error) {

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, errors.New("given tree type is not a struct")
	}

	tr := new(tree)
	tr.Type = t
	tr.fields = make(map[string]*field)

	err := iteratePublicFields(t, func(idx int, field *reflect.StructField) error {
		fx, err := buildField(tag, idx, field)
		if err != nil {
			return err
		}
		tr.fields[fx.Identifier.Name] = fx
		return nil
	})
	if err != nil {
		return nil, err
	}

	return tr, nil
}

func iteratePublicFields(t reflect.Type, f func(int, *reflect.StructField) error) error {

	for i := t.NumField() - 1; i >= 0; i-- {
		st := t.Field(i)
		if err := f(i, &st); err != nil {
			return err
		}
	}

	return nil
}

func buildField(tag string, idx int, sf *reflect.StructField) (f *field, err error) {
	f = new(field)
	f.Index = idx
	f.Type = sf.Type

	if tag == "" {
		// add the non-tag version
		f.Identifier.Name = sf.Name
	} else {
		// add the tag identifier
		val := sf.Tag.Get(tag)
		item := strings.Split(val, ",")[0]
		f.Identifier.Name = item
	}

	switch sf.Type.Kind() {
	case reflect.Ptr:
		if sf.Type.Elem().Kind() == reflect.Struct {
			f.tree, err = buildTree(tag, sf.Type)
			if err != nil {
				return
			}
		}
	case reflect.Struct:
		f.tree, err = buildTree(tag, sf.Type)
	default:
	}

	return
}
