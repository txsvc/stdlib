package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/johngb/langreg"
)

const (
	// AssertionWarning indicates a potential issue
	AssertionWarning = 0
	// AssertionError indicates an error in the validation
	AssertionError = 1

	MsgStringMismatch      = "expected '%s', found '%s'"
	MsgNonEmptyString      = "expected non empty string '%s'"
	MsgNotNil              = "expected non-nil attribute '%s'"
	MsgNonZero             = "expected non-zero value '%s'"
	MsgInvalidLanguageCode = "invalid language code '%s'"
	MsgInvalidTimestamp    = "invalid timestamp '%d'"
	MsgNonEmptyMap         = "expected none empty map '%s'"
	MsgExpectedKey         = "expected key '%s' in map '%s'"
	MsgEmptyReport         = "no errors/warnings"
)

type (
	// Assertion is used to collect validation information
	Assertion struct {
		Type int    // 0 == warning, 1 == error
		Txt  string // description of the problem
		Err  error
		ctx  string // context describes where the error or warning happened

	}

	// Validator collects assertions
	Validator struct {
		Issues   []*Assertion
		Errors   int
		Warnings int
		ctxStack []string
	}
)

// NewValidator initializes and returns a new Validator
func NewValidator() *Validator {
	v := Validator{
		Issues:   make([]*Assertion, 0),
		Errors:   0,
		Warnings: 0,
	}
	v.ctxStack = make([]string, 0)

	return &v
}

// AddError adds an error assertion
func (v *Validator) AddError(txt string) {
	v.Issues = append(v.Issues, &Assertion{Type: AssertionError, Txt: txt, ctx: v.Context()})
	v.Errors++
}

// AddWarning adds an warning assertion
func (v *Validator) AddWarning(txt string) {
	v.Issues = append(v.Issues, &Assertion{Type: AssertionWarning, Txt: txt, ctx: v.Context()})
	v.Warnings++
}

// StringEquals verifies a string
func (v *Validator) StringEquals(src, expected string) {
	if len(src) != len(expected) {
		v.AddError(fmt.Sprintf(MsgStringMismatch, expected, src))
		return
	}

	if src != expected {
		v.AddError(fmt.Sprintf(MsgStringMismatch, expected, src))
	}
}

// StringNotEmpty verifies a string is not empty
func (v *Validator) StringNotEmpty(src, hint string) {
	if len(src) == 0 {
		v.AddError(fmt.Sprintf(MsgNonEmptyString, hint))
	}
}

// NotNil verifies that an attribute is not nil
func (v *Validator) NotNil(src interface{}, hint string) {
	if src == nil || (reflect.ValueOf(src).Kind() == reflect.Ptr && reflect.ValueOf(src).IsNil()) {
		v.AddError(fmt.Sprintf(MsgNotNil, hint))
	}
}

// NonZero verifies that a map is not empty
func (v *Validator) NonZero(src int, hint string) {
	if src == 0 {
		v.AddError(fmt.Sprintf(MsgNonZero, hint))
	}
}

// ISO639 verifies that src complies with ISO 639-1
func (v *Validator) ISO639(src string) {
	lang := src
	if !strings.Contains(src, "_") {
		lang = src + "_" + strings.ToUpper(src)
	}
	if !langreg.IsValidLangRegCode(lang) {
		v.AddError(fmt.Sprintf(MsgInvalidLanguageCode, src))
	}
}

// RFC1123Z verifies that src complies with RFC 1123-Z
func (v *Validator) RFC1123Z(src string) {
	_, err := time.Parse(time.RFC1123Z, src)
	if err != nil {
		v.AddError(err.Error())
	}
}

// Timestamp validates that src is a valid UNIX timestamp
func (v *Validator) Timestamp(src string) {
	ts, err := strconv.Atoi(src)
	if err != nil {
		v.AddError(err.Error())
	}
	if ts < 0 {
		v.AddError(fmt.Sprintf(MsgInvalidTimestamp, ts))
	}
}

// MapNotEmpty verifies that a map is not empty
func (v *Validator) MapNotEmpty(src map[string]string, hint string) {
	if len(src) == 0 {
		v.AddError(fmt.Sprintf(MsgNonEmptyMap, hint))
	}
}

// MapContains verifies that a map contains key
func (v *Validator) MapContains(src map[string]string, key, hint string) {
	if len(src) == 0 {
		v.AddError(fmt.Sprintf(MsgNonEmptyMap, hint))
		return
	}
	if _, ok := src[key]; !ok {
		v.AddError(fmt.Sprintf(MsgExpectedKey, key, hint))
	}
}

// IsValid returns true if NError == 0. Warnings are ignored
func (v *Validator) IsValid() bool {
	return v.Errors == 0
}

// IsClean returns true if NError == 0 AND NWarnings == 0
func (v *Validator) IsClean() bool {
	return v.Errors == 0 && v.Warnings == 0
}

// NErrors returns the number of erros
func (v *Validator) NErrors() int {
	return v.Errors
}

// NWarnings returns the number of warnings
func (v *Validator) NWarnings() int {
	return v.Warnings
}

// AsError returns an error if NError > 0, nil otherwise
func (v *Validator) AsError() error {
	if v.Errors == 0 {
		return nil
	}
	return fmt.Errorf(v.Error())
}

// Error returns an error text
func (v *Validator) Error() string {
	return v.Report()
}

// Report returns a description of all issues
func (v *Validator) Report() string {
	if v.Errors == 0 && v.Warnings == 0 {
		return MsgEmptyReport
	}
	r := "\n"
	for i, issue := range v.Issues {
		r = r + issue.ToStringWithIndex(i+1)
	}
	return r
}

func (a *Assertion) ToString() string {
	return fmt.Sprintf("%s(%s): %s\n", a.ctx, a.TypeAsString(), a.Txt)
}

func (a *Assertion) ToStringWithIndex(i int) string {
	return fmt.Sprintf("%s-%02d (%s): %s\n", a.TypeAsString(), i, a.ctx, a.Txt)
}

func (a *Assertion) TypeAsString() string {
	if a.Type == AssertionError {
		return "error"
	}
	return "warning"
}
