package stdlib

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	testCases := []struct {
		name     string
		level    string
		expected zerolog.Level
	}{
		{"trace level", "trace", zerolog.TraceLevel},
		{"debug level", "debug", zerolog.DebugLevel},
		{"info level", "info", zerolog.InfoLevel},
		{"warn level", "warn", zerolog.WarnLevel},
		{"error level", "error", zerolog.ErrorLevel},
		{"invalid level", "invalid", zerolog.Disabled},
		{"empty level", "", zerolog.Disabled},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.level != "" {
				os.Setenv(LOG_LEVEL, tc.level)
				defer os.Unsetenv(LOG_LEVEL)
			}
			SetLogLevel()
			assert.Equal(t, tc.expected, zerolog.GlobalLevel())
		})
	}
}

func TestPlaceholder(t *testing.T) {
	// Test that Placeholder is empty struct
	assert.Equal(t, struct{}{}, Placeholder)

	// Test AnyType can hold different types
	var anyVal AnyType
	anyVal = 42
	assert.Equal(t, 42, anyVal)
	anyVal = "test"
	assert.Equal(t, "test", anyVal)
	anyVal = true
	assert.Equal(t, true, anyVal)
}
