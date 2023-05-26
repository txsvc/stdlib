package stdlib

import (
	"os"
	"strconv"
	"strings"
)

// GetString returns the environment variable ENV['env'] or def if 'env' is not set.
// Note: def is only returned if the 'env' is not set, i.e. an EMPTY 'env' is still returned !
func GetString(env, def string) string {
	e, ok := os.LookupEnv(env)
	if ok {
		return e
	}
	return def
}

// GetInt returns the environment variable ENV['env'] as an int64 or def if 'env' is not set.
func GetInt(env string, def int64) int64 {
	e, ok := os.LookupEnv(env)
	if ok {
		v, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return def
		}
		return v
	}
	return def
}

// GetBool returns the environment variable ENV['env'] as boolean or def if 'env' is not set.
func GetBool(env string, def bool) bool {
	e, ok := os.LookupEnv(env)
	if !ok {
		return def
	}

	e = strings.ToLower(e)
	if e == "true" || e == "yes" || e == "1" {
		return true
	}
	return false
}

// Exists verifies that the environment variable 'env' is defined and returns a non-empty value.
func Exists(env string) bool {
	e, ok := os.LookupEnv(env)
	if !ok || e == "" {
		return false
	}
	return true
}
