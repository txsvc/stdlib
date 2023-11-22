package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	assert.NotNil(t, v)
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestStringEquals(t *testing.T) {
	v := NewValidator()

	v.StringEquals("a", "a")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.StringEquals("a", "b")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestStringNotEmpty(t *testing.T) {
	v := NewValidator()

	v.StringNotEmpty("a", "hint")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.StringNotEmpty("", "hint")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestNotNil(t *testing.T) {
	v1 := NewValidator()
	v1.NotNil("a", "hint")
	assert.Equal(t, 0, v1.Errors)
	assert.Equal(t, 0, v1.Warnings)

	v2 := NewValidator()
	v2.NotNil("", "hint")
	assert.Equal(t, 0, v2.Errors)
	assert.Equal(t, 0, v2.Warnings)

	var s1 *Assertion
	v3 := NewValidator()
	v3.NotNil(s1, "hint")
	assert.Equal(t, 1, v3.Errors)
	assert.Equal(t, 0, v3.Warnings)

	v4 := NewValidator()
	s2 := new(Assertion)
	v4.NotNil(s2, "hint")
	assert.Equal(t, 0, v4.Errors)
	assert.Equal(t, 0, v4.Warnings)
}

func TestNonZero(t *testing.T) {
	v := NewValidator()

	v.NonZero(42, "hint")
	assert.Equal(t, 0, v.Errors)
	v.NonZero(-66, "hint")
	assert.Equal(t, 0, v.Errors)
	v.NonZero(0, "hint")
	assert.Equal(t, 1, v.Errors)
}

func TestISO639(t *testing.T) {
	v := NewValidator()

	v.ISO639("en_US")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)
	v.ISO639("en_us")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestRFC1123Z(t *testing.T) {
	v := NewValidator()

	v.RFC1123Z("Fri, 04 Jun 2021 06:57:26 +0000")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.RFC1123Z("04 Jun 2021 06:57:26 +0000")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.RFC1123Z("Fri, 04 Jun 2021 06:57:26")
	assert.Equal(t, 2, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestTimestamp(t *testing.T) {
	v := NewValidator()

	v.Timestamp("1")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.Timestamp("-1")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.Timestamp("sdfdfd")
	assert.Equal(t, 2, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestMapNotEmpty(t *testing.T) {

	m := make(map[string]string)

	v1 := NewValidator()
	v1.MapNotEmpty(m, "hint")
	assert.Equal(t, 1, v1.Errors)
	assert.Equal(t, 0, v1.Warnings)

	v2 := NewValidator()
	m["foo"] = "bar"
	v2.MapNotEmpty(m, "hint")
	assert.Equal(t, 0, v2.Errors)
	assert.Equal(t, 0, v2.Warnings)
}

func TestMapContains(t *testing.T) {
	v := NewValidator()
	m := make(map[string]string)
	m["foo"] = "bar"

	v.MapContains(m, "foo", "hint")
	assert.Equal(t, 0, v.Errors)
	assert.Equal(t, 0, v.Warnings)

	v.MapContains(m, "bar", "hint")
	assert.Equal(t, 1, v.Errors)
	assert.Equal(t, 0, v.Warnings)
}

func TestReport(t *testing.T) {
	m := make(map[string]string)

	v1 := NewValidator()
	v1.MapNotEmpty(m, "hint 1")
	v1.AddWarning("warning 1")

	fmt.Println(v1.AsError())
}

func TestSaveAndResoreContext(t *testing.T) {
	v1 := NewValidator()
	assert.Equal(t, 0, len(v1.ctxStack))
	assert.Equal(t, "root", v1.Context())

	v1.SaveContext("ctx1")
	assert.Equal(t, 1, len(v1.ctxStack))
	assert.Equal(t, "ctx1", v1.Context())

	v1.SaveContext("ctx2")
	assert.Equal(t, 2, len(v1.ctxStack))
	assert.Equal(t, "ctx2", v1.Context())

	v1.RestoreContext()
	assert.Equal(t, 1, len(v1.ctxStack))
	assert.Equal(t, "ctx1", v1.Context())

	v1.RestoreContext()
	assert.Equal(t, 0, len(v1.ctxStack))
	assert.Equal(t, "root", v1.Context())
}
