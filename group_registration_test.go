package typetree

import (
	"errors"
	"fmt"
	"testing"
)

type testType struct {
}

var groupRegistrationTests = []groupRegistrationTest{
	{Name: "", Input: nil, Expected: errors.New("Type argument is nil")},
	{Name: "name1", Input: nil, Expected: errors.New("Type argument is nil")},

	{Name: "", Input: testType{}, Expected: errors.New("Name is invalid")},
	{Name: "name1", Input: testType{}, Expected: nil},

	{Name: "name1", Input: 12, Expected: errors.New("Type 'int' is not a struct")},
}

type groupRegistrationTest struct {
	Name  string
	Tag   string
	Input interface{}

	Expected error
}

func (g *groupRegistrationTest) FailureString(i ...interface{}) string {
	return fmt.Sprintf("group.Register('%v', '%v', '%v') => %v, expected %v",
		g.Name, g.Tag, g.Input, i, g.Expected,
	)
}

func (g *groupRegistrationTest) Test(err error) bool {
	return err == g.Expected || (err != nil && g.Expected != nil && err.Error() == g.Expected.Error())
}

func TestGroupRegistration(t *testing.T) {
	for _, gtest := range groupRegistrationTests {
		var group Group
		err := group.Register(gtest.Name, gtest.Tag, gtest.Input)

		if !gtest.Test(err) {
			t.Errorf(gtest.FailureString(err))
		}
	}
}
