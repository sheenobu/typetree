package typetree

import (
	"errors"
	"fmt"
	"testing"
)

type groupNewTest struct {
	Group *Group
	Name  string

	ExpectedOutput interface{}
	ExpectedError  error
}

func (g *groupNewTest) FailureString(i ...interface{}) string {
	return fmt.Sprintf("group.New(%v) => (%v, %v), expected (%v, %v)",
		g.Name, i[0], i[1], g.ExpectedOutput, g.ExpectedError)
}

func (g *groupNewTest) Test(out interface{}, err error) bool {
	success := (err == g.ExpectedError || (err != nil && g.ExpectedError != nil && err.Error() == g.ExpectedError.Error()))
	success = success && (out == g.ExpectedOutput)

	return success
}

var emptyGroup = &Group{}
var nonEmptyGroup = &Group{}

type gx struct {
}

func init() {
	nonEmptyGroup.Register("hello", "", &gx{})
}

var groupNewTests = []groupNewTest{
	{Group: emptyGroup, Name: "", ExpectedOutput: nil, ExpectedError: errors.New("Name is invalid")},
	{Group: emptyGroup, Name: "hello", ExpectedOutput: nil, ExpectedError: errors.New("Given type 'hello' not found")},

	{Group: nonEmptyGroup, Name: "", ExpectedOutput: nil, ExpectedError: errors.New("Name is invalid")},
	{Group: nonEmptyGroup, Name: "hello", ExpectedOutput: &gx{}, ExpectedError: nil},

	{Group: nonEmptyGroup, Name: "world", ExpectedOutput: nil, ExpectedError: errors.New("Given type 'world' not found")},
}

func TestGroupNew(t *testing.T) {
	for _, gtest := range groupNewTests {
		out, err := gtest.Group.New(gtest.Name)

		var v interface{}
		if out != nil {
			v = out.Interface()
		}

		if !gtest.Test(v, err) {
			t.Errorf(gtest.FailureString(v, err))
		}
	}
}
