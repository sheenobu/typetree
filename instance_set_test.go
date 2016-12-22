package typetree

import "testing"

type itt struct {
	Hello string `json:"h2"`
	World string `json:"w2"`
}

type ittComplex struct {
	World string `json:"w2"`
	I     *itt   `json:"itt"`
	I2    itt    `json:"itt2"`
}

type ittComplex2 struct {
	IC ittComplex `json:"x"`
}

func init() {
	nonEmptyGroup.Register("instanceTestType", "", &itt{})
	nonEmptyGroup.Register("instanceTestTypeWithTag", "json", &itt{})
	nonEmptyGroup.Register("instanceTestTypeComplexWithTag", "json", &ittComplex{})
	nonEmptyGroup.Register("ittComplex2", "json", &ittComplex2{})

}

func TestInstanceSetSimple(t *testing.T) {

	instance, err := nonEmptyGroup.New("instanceTestType")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	var i3 = instance.Interface().(*itt)

	err = instance.Set(Keys("Hello"), "world")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	var ix = instance.Interface().(*itt)
	if ix.Hello != "world" {
		t.Errorf("mismatched")
	}

	if i3.Hello != "world" {
		t.Errorf("mismatched")
	}

}

func TestInstanceSetSimpleTag(t *testing.T) {

	instance, err := nonEmptyGroup.New("instanceTestTypeWithTag")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	var i3 = instance.Interface().(*itt)

	err = instance.Set(Keys("h2"), "world")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	var ix = instance.Interface().(*itt)
	if ix.Hello != "world" {
		t.Errorf("mismatched")
	}

	if i3.Hello != "world" {
		t.Errorf("mismatched")
	}

}

func TestInstanceSetComplexTag(t *testing.T) {

	instance, err := nonEmptyGroup.New("instanceTestTypeComplexWithTag")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	err = instance.Set(Keys("h2"), "world")
	if err == nil {
		t.Errorf("Expected error")
	}

	err = instance.Set(Keys("w2"), "x")
	if err != nil {
		t.Errorf("Err: %v", err)
	}

	err = instance.Set(Keys("itt", "h2"), "z")
	if err != nil {
		t.Errorf("Err: %v", err)
	}

	err = instance.Set(Keys("itt2", "h2"), "z2")
	if err != nil {
		t.Errorf("Err: %v", err)
	}

	results := instance.Interface().(*ittComplex)
	if results.World != "x" {
		t.Errorf("Mismatch world")
	}

	if results.I == nil {
		t.Errorf("I is nil")
	}

	if results.I.Hello != "z" {
		t.Errorf("itt/h2 mismatch")
	}

	if results.I.World != "" {
		t.Errorf("itt/w2 mismatch")
	}

	if results.I2.Hello != "z2" {
		t.Errorf("itt2/z2 mismatch")
	}
}

func TestInstanceSetComplexTag2(t *testing.T) {

	instance, err := nonEmptyGroup.New("ittComplex2")
	if err != nil {
		t.Errorf("err: %v", err)
	}

	err = instance.Set(Keys("h2"), "world")
	if err == nil {
		t.Errorf("Expected error")
	}

	err = instance.Set(Keys("x", "itt2", "h2"), "z")
	if err != nil {
		t.Errorf("Err: %v", err)
	}

	results := instance.Interface().(*ittComplex2)
	if results.IC.I2.Hello != "z" {
		t.Errorf("Mismatch x/itt2/h2/")
	}

}
