package typetree

import (
	"reflect"

	"github.com/pkg/errors"
)

// An Instance represents a struct and the series of
// operations you can perform on the struct to both inspect
// and update the struct and its children.
type Instance interface {

	// Interface returns the pointer value of the object.
	Interface() interface{}

	// Set sets the given value on the structure, if it exists.
	Set(key Key, value interface{}) error

	// Get sets the given value on the structure, if it exists.
	Get(key Key) (Value, error)
}

// implementation of Instance that uses the internal tree structures
type treeInstance struct {
	tree *tree
	val  reflect.Value
}

func (ti *treeInstance) Interface() interface{} {
	return ti.val.Interface()
}

func (ti *treeInstance) Set(key Key, value interface{}) (err error) {
	var fx = ti.tree.fields
	var curPtr = ti.val

	// if there are any struct types that aren't wrapped as pointers,
	// we must first create new instance of that struct type, operate on that,
	// then perform the copy.

	var copies []func() error
	defer func() {
		if err == nil {

			// we iterate backwards since the tree is built forwards.
			// Test 'TestInstanceSetComplexTag2' confirms this.

			for i := len(copies) - 1; i >= 0; i-- {
				x := copies[i]
				if err = x(); err != nil {
					return
				}
			}
		}
	}()

	for _, k := range key {
		if child, ok := fx[k]; ok {
			if child.tree == nil {
				// simple type
				err = errors.Wrap(child.SetValue(curPtr.Interface(), value), "Failed to set simple value")
				return
			}

			// complex type
			var i interface{}
			i, err = child.GetValue(curPtr.Interface())
			if err != nil {
				return
			}

			if i == nil {
				i = reflect.New(child.Type.Elem()).Interface()
				if err := child.SetValue(curPtr.Interface(), i); err != nil {
					return err
				}
			} else if reflect.TypeOf(i).Kind() == reflect.Struct {
				nextPtr := reflect.New(reflect.TypeOf(i))

				copies = append(copies, func(c *field, iface interface{}, nptr *reflect.Value) func() error {
					return func() error {
						return c.SetValue(iface, nptr.Interface())
					}
				}(child, curPtr.Interface(), &nextPtr))

				curPtr = nextPtr
				fx = child.tree.fields

				continue
			}

			fx = child.tree.fields
			curPtr = reflect.ValueOf(i)
		}
	}

	return errors.Errorf("Could not find field '%v'", key)
}

func (ti *treeInstance) Get(key Key) (val Value, err error) {
	var fx = ti.tree.fields
	var curPtr = ti.val

	for _, k := range key {
		if child, ok := fx[k]; ok {
			var i interface{}
			if child.tree == nil {
				// simple type
				i, err = child.GetValue(curPtr.Interface())
				if err != nil {
					return
				}

				val = &valueS{
					t: child.Type,
					v: reflect.ValueOf(i),
				}
				return
			}

			// complex type
			i, err = child.GetValue(curPtr.Interface())
			if err != nil {
				return
			}

			if i == nil {
				err = errors.Errorf("Could not find field '%v'", key)
				return
			}

			fx = child.tree.fields
			curPtr = reflect.ValueOf(i)
		}
	}

	err = errors.Errorf("Could not find field '%v'", key)
	return
}
