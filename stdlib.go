package stdlib

import (
	"math/rand"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	LOG_LEVEL = "LOG_LEVEL"
)

// Placeholder is a placeholder object that can be used globally.
var Placeholder PlaceholderType
var src = rand.NewSource(time.Now().UnixNano())

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

// Seed sets the seed for the default random source.
// This affects all random string generation functions that don't use crypto/rand.
func Seed(seed int64) {
	src.Seed(seed)
}
