package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	Validate1 struct {
		attr1 string
		attr2 int
	}

	Validate2 struct {
		attr1 string
		attr2 *Validate1
	}

	Validate3 struct {
		attr1 string
		attr2 *Validate1
	}
)

func (m *Validate1) Validate(root string, v *Validator) *Validator {
	v.SaveContext(root)
	defer v.RestoreContext()

	v.StringNotEmpty(m.attr1, "attr1")
	v.NonZero(m.attr2, "attr2")
	if len(m.attr1) < 10 {
		v.AddWarning("attr1 should be longer")
	}

	return v
}

func (m *Validate2) Validate(root string, v *Validator) *Validator {
	v.SaveContext(root)
	defer v.RestoreContext()

	v.StringNotEmpty(m.attr1, "attr1")
	return m.attr2.Validate(root, v)
}

func (m *Validate3) Validate(root string, v *Validator) *Validator {
	v.SaveContext(root)
	defer v.RestoreContext()

	v.StringNotEmpty(m.attr1, "attr1")
	return m.attr2.Validate("Validate3.attr2", v)
}

func TestSimpleSuccess(t *testing.T) {

	vs1 := Validate1{
		attr1: "some string",
		attr2: 42,
	}

	v := NewValidator()
	v = vs1.Validate("TestSimpleSuccess", v)

	assert.NotNil(t, v)
	assert.Equal(t, 0, v.Errors)
}

func TestSimpleFail(t *testing.T) {

	vs1 := Validate1{
		attr1: "some string",
	}

	v := NewValidator()
	v = vs1.Validate("TestSimpleFail", v)

	assert.NotNil(t, v)
	assert.Equal(t, 0, v.Warnings)
	assert.Equal(t, 1, v.Errors)
}

func TestSimpleWarning(t *testing.T) {

	vs1 := Validate1{
		attr1: "string",
		attr2: 1,
	}

	v := NewValidator()
	v = vs1.Validate("TestSimpleWarning", v)
	assert.NotNil(t, v)

	fmt.Println(v.Report())

	assert.Equal(t, 1, v.Warnings)
	assert.Equal(t, 0, v.Errors)
}

func TestSimpleNested(t *testing.T) {

	vs2 := Validate1{
		attr1: "some string",
		attr2: 42,
	}

	vs1 := Validate2{
		attr1: "some string",
		attr2: &vs2,
	}

	v := NewValidator()
	v = vs1.Validate("TestSimpleNested", v)

	assert.NotNil(t, v)
	assert.Equal(t, 0, v.Errors)
}

func TestNestedWithError(t *testing.T) {

	vs2 := Validate1{
		attr1: "some string",
		attr2: 0,
	}

	vs1 := Validate3{
		attr2: &vs2,
	}

	v := NewValidator()
	v = vs1.Validate("TestNestedWithError", v)

	assert.NotNil(t, v)
	assert.Equal(t, 2, v.Errors)

	fmt.Println(v.Report())
}
