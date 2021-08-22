package provider

import (
	"context"
)

const (
	TypeLogger ProviderType = iota
	TypeErrorReporter
	TypeMetrics

	LevelInfo Severity = iota
	LevelWarn
	LevelError
	LevelDebug
)

type (
	ProviderType int

	InstanceProviderFunc func() interface{}

	ProviderConfig struct {
		ID   string
		Type ProviderType
		Impl InstanceProviderFunc
	}

	GenericProvider interface {
		Close() error
	}
)

type (
	Severity int

	ErrorReportingProvider interface {
		ReportError(error)
	}

	MetricsProvider interface {
		Meter(ctx context.Context, metric string, args ...string)
	}

	LoggingProvider interface {
		Log(string, ...string)
		LogWithLevel(Severity, string, ...string)
	}
)

// Returns the name of a provider type
func (l ProviderType) String() string {
	switch l {
	case TypeLogger:
		return "LOGGER"
	case TypeErrorReporter:
		return "ERROR_REPORTER"
	case TypeMetrics:
		return "METRICS"
	default:
		panic("unsupported")
	}
}

// WithProvider returns a populated ProviderConfig struct.
func WithProvider(ID string, providerType ProviderType, impl InstanceProviderFunc) ProviderConfig {
	return ProviderConfig{
		ID:   ID,
		Type: providerType,
		Impl: impl,
	}
}
