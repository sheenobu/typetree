package typetree

import "testing"

type arrayTestStruct struct {
	Items []string `json:"items"`
}

func init() {
	nonEmptyGroup.Register("arrayTestStruct", "json", &arrayTestStruct{})
}

func TestArrayAppend(t *testing.T) {

	instance, err := nonEmptyGroup.New("arrayTestStruct")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	val, err := instance.Get(Keys("items"))
	if err != nil {
		t.Errorf("err: %v", err)
	}

	err = val.Append("hello")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	s := instance.Interface().(*arrayTestStruct)
	if len(s.Items) != 1 {
		t.Errorf("item length should be 1")
	}
}
