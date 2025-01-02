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
	if log_level == "trace" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if log_level == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else if log_level == "info" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if log_level == "warn" {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if log_level == "error" {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
