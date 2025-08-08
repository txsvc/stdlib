package stdlib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	e := GetString("PATH", "nix")

	assert.NotEmpty(t, e)
	assert.NotEqual(t, "nix", e)
}

func TestGetStringDefault(t *testing.T) {
	e := GetString("FOO_BAR", "nix")

	assert.NotEmpty(t, e)
	assert.Equal(t, "nix", e)
}

func TestGetIntDefault(t *testing.T) {
	e := GetInt("FOO_BAR", 42)
	assert.Equal(t, int64(42), e)
}

func TestGetIntNotInt(t *testing.T) {
	e := GetInt("PATH", 66)
	assert.Equal(t, int64(66), e)
}

func TestAssert(t *testing.T) {
	assert.True(t, Exists("PATH"))
	assert.False(t, Exists("FOO_BAR"))
}

func TestGetBool(t *testing.T) {
	testCases := []struct {
		name     string
		envKey   string
		envValue string
		def      bool
		expected bool
	}{
		{"true value", "TEST_BOOL_1", "true", false, true},
		{"yes value", "TEST_BOOL_2", "yes", false, true},
		{"1 value", "TEST_BOOL_3", "1", false, true},
		{"false value", "TEST_BOOL_4", "false", true, false},
		{"invalid value", "TEST_BOOL_5", "invalid", true, false},
		{"empty value", "TEST_BOOL_6", "", false, false},
		{"default true", "TEST_BOOL_NOT_SET", "", true, true},
		{"default false", "TEST_BOOL_NOT_SET", "", false, false},
		{"uppercase TRUE", "TEST_BOOL_7", "TRUE", false, true},
		{"uppercase YES", "TEST_BOOL_8", "YES", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				_ = os.Setenv(tc.envKey, tc.envValue)
				defer func() { _ = os.Unsetenv(tc.envKey) }()
			}

			result := GetBool(tc.envKey, tc.def)
			assert.Equal(t, tc.expected, result)
		})
	}
}
