package stdlib

import (
	"strings"

	"github.com/rs/zerolog"
)

const (
	LOG_LEVEL = "LOG_LEVEL"
)

// Placeholder is a placeholder object that can be used globally.
var Placeholder PlaceholderType

type (
	// AnyType can be used to hold any type.
	AnyType = interface{}
	// PlaceholderType represents a placeholder type.
	PlaceholderType = struct{}
)

// SetLogLevel configures the global logging level based on the LOG_LEVEL environment variable.
// Supported values are: "trace", "debug", "info", "warn", "error".
// If the value is not recognized or not set, logging will be disabled.
func SetLogLevel() {
	// setup logging
	log_level := strings.ToLower(GetString(LOG_LEVEL, ""))
	switch log_level {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
