package typetree

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type fieldSetTest struct {
	Subject     field
	Value       interface{}
	Destination interface{}

	ExpectedX     x
	ExpectedError error
}

func (g *fieldSetTest) FailureString(i ...interface{}) string {
	return fmt.Sprintf("field{%v}.Set(%v) => %v; %v, expected %v; %v",
		g.Subject, g.Value, i[0], i[1], g.ExpectedError, g.ExpectedX)
}

func (g *fieldSetTest) Test(i interface{}, err error) bool {

	var ex = i.(*x)

	success := (err == g.ExpectedError || (err != nil && g.ExpectedError != nil && err.Error() == g.ExpectedError.Error()))
	success = success && g.ExpectedX.Field1 == ex.Field1
	success = success && g.ExpectedX.Field2 == ex.Field2
	success = success && reflect.DeepEqual(g.ExpectedX.Field3, ex.Field3)
	success = success && reflect.DeepEqual(g.ExpectedX.Field4, ex.Field4)

	return success
}

type x struct {
	Field1 string
	Field2 int
	Field3 *y
	Field4 y
	Field5 []string
}

type y struct {
	X int32
}

var field1 = field{Index: 0, Identifier: identifier{Name: "Field1"}, Type: reflect.TypeOf("")}
var field2 = field{Index: 1, Identifier: identifier{Name: "Field2"}, Type: reflect.TypeOf(0)}
var field3 = field{Index: 2, Identifier: identifier{Name: "Field3"}, Type: reflect.TypeOf(&y{})}
var field4 = field{Index: 3, Identifier: identifier{Name: "Field4"}, Type: reflect.TypeOf(y{})}
var field5 = field{Index: 4, Identifier: identifier{Name: "Field5"}, Type: reflect.TypeOf([]string{})}

var fieldSetTests = []fieldSetTest{
	{Subject: field1, Value: "hello", Destination: &x{}, ExpectedX: x{Field1: "hello"}, ExpectedError: nil},
	{Subject: field2, Value: 12, Destination: &x{}, ExpectedX: x{Field2: 12}, ExpectedError: nil},

	{Subject: field1, Value: 12, Destination: &x{}, ExpectedX: x{}, ExpectedError: errors.New("Value '12' is not assignable to field")},
	{Subject: field2, Value: "hello", Destination: &x{}, ExpectedX: x{}, ExpectedError: errors.New("Value 'hello' is not assignable to field")},

	{Subject: field3, Value: &y{X: 12}, Destination: &x{}, ExpectedX: x{Field3: &y{X: 12}}, ExpectedError: nil},
	{Subject: field3, Value: y{X: 12}, Destination: &x{}, ExpectedX: x{Field3: &y{X: 12}}, ExpectedError: nil},

	{Subject: field4, Value: &y{X: 12}, Destination: &x{}, ExpectedX: x{Field4: y{X: 12}}, ExpectedError: nil},
	{Subject: field4, Value: y{X: 12}, Destination: &x{}, ExpectedX: x{Field4: y{X: 12}}, ExpectedError: nil},

	{Subject: field5, Value: []string{"hello"}, Destination: &x{}, ExpectedX: x{Field5: []string{"hello"}}, ExpectedError: nil},
}

func TestFieldSet(t *testing.T) {
	for _, ftest := range fieldSetTests {
		err := ftest.Subject.SetValue(ftest.Destination, ftest.Value)
		if !ftest.Test(ftest.Destination, err) {
			t.Errorf(ftest.FailureString(err, ftest.Destination))
		}
	}
}
