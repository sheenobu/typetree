package typetree

import "testing"

func TestInstanceGetSimple(t *testing.T) {

	instance, err := nonEmptyGroup.New("instanceTestType")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	var i3 = instance.Interface().(*itt)
	i3.Hello = "world"

	val, err := instance.Get(Keys("Hello"))
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if val.Interface().(string) != "world" {
		t.Errorf("Failed to get hello key")
	}

}

func TestInstanceGetSimpleTag(t *testing.T) {

	instance, err := nonEmptyGroup.New("instanceTestTypeWithTag")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	var i3 = instance.Interface().(*itt)
	i3.Hello = "world"

	val, err := instance.Get(Keys("h2"))
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if val.Interface().(string) != "world" {
		t.Errorf("Failed to get h2 key")
	}

}

func TestInstanceGetComplexTag(t *testing.T) {

	instance, err := nonEmptyGroup.New("instanceTestTypeComplexWithTag")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	var i3 = instance.Interface().(*ittComplex)
	i3.World = "x"
	i3.I = &itt{Hello: "z"}
	i3.I2.Hello = "z2"

	_, err = instance.Get(Keys("h2"))
	if err == nil {
		t.Errorf("Expected error")
	}

	val, err := instance.Get(Keys("w2"))
	if err != nil {
		t.Errorf("Err: %v", err)
	}
	if val.Interface().(string) != "x" {
		t.Errorf("Failed to get h2 key")
	}

	val, err = instance.Get(Keys("itt", "h2"))
	if err != nil {
		t.Errorf("Err: %v", err)
	}
	if val.Interface().(string) != "z" {
		t.Errorf("Failed to get itt/h2 key")
	}

	val, err = instance.Get(Keys("itt2", "h2"))
	if err != nil {
		t.Errorf("Err: %v", err)
	}
	if val.Interface().(string) != "z2" {
		t.Errorf("Failed to get itt2/h2 key")
	}

}
