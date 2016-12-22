package typetree

import (
	"fmt"
	"reflect"
	"testing"
)

type fieldGetTest struct {
	Subject field
	Source  interface{}

	ExpectedOut   interface{}
	ExpectedError error
}

func (g *fieldGetTest) FailureString(i ...interface{}) string {
	return fmt.Sprintf("field{%v}.Get(%v) => %v, %v, expected %v, %v",
		g.Subject, g.Source, i[0], i[1], g.ExpectedOut, g.ExpectedError)
}

func (g *fieldGetTest) Test(i interface{}, err error) bool {

	success := (err == g.ExpectedError || (err != nil && g.ExpectedError != nil && err.Error() == g.ExpectedError.Error()))
	success = success && reflect.DeepEqual(g.ExpectedOut, i)

	return success
}

var fieldGetTests = []fieldGetTest{
	{Subject: field1, Source: &x{Field1: "hello"}, ExpectedOut: "hello", ExpectedError: nil},
	{Subject: field1, Source: x{Field1: "hello"}, ExpectedOut: "hello", ExpectedError: nil},

	{Subject: field2, Source: &x{Field2: 12}, ExpectedOut: 12, ExpectedError: nil},
	{Subject: field2, Source: x{Field2: 12}, ExpectedOut: 12, ExpectedError: nil},

	{Subject: field3, Source: &x{Field3: &y{X: 12}}, ExpectedOut: &y{X: 12}, ExpectedError: nil},
	{Subject: field3, Source: x{Field3: &y{X: 12}}, ExpectedOut: &y{X: 12}, ExpectedError: nil},

	{Subject: field4, Source: &x{Field4: y{X: 14}}, ExpectedOut: y{X: 14}, ExpectedError: nil},
	{Subject: field4, Source: x{Field4: y{X: 14}}, ExpectedOut: y{X: 14}, ExpectedError: nil},

	{Subject: field5, Source: &x{Field5: []string{"hello"}}, ExpectedOut: []string{"hello"}, ExpectedError: nil},
}

func TestFieldGet(t *testing.T) {
	for _, ftest := range fieldGetTests {
		val, err := ftest.Subject.GetValue(ftest.Source)
		if !ftest.Test(val, err) {
			t.Errorf(ftest.FailureString(val, err))
		}
	}
}
