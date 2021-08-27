package observer

import (
	"context"

	"github.com/txsvc/stdlib/pkg/provider"
)

const (
	LevelInfo Severity = iota
	LevelWarn
	LevelError
	LevelDebug

	TypeLogger provider.ProviderType = iota
	TypeErrorReporter
	TypeMetrics
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

	Observer struct {
		Provider *provider.Provider
	}
)

var (
	observer *Observer
)

func (obs *Observer) Log(msg string, keyValuePairs ...string) {
	imp, found := obs.Provider.Find(TypeLogger)
	if !found {
		return
	}
	imp.(LoggingProvider).Log(msg, keyValuePairs...)
}

func (obs *Observer) LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
	imp, found := obs.Provider.Find(TypeLogger)
	if !found {
		return
	}
	imp.(LoggingProvider).LogWithLevel(lvl, msg, keyValuePairs...)
}

func (obs *Observer) Meter(ctx context.Context, metric string, args ...string) {
	imp, found := obs.Provider.Find(TypeMetrics)
	if !found {
		return
	}
	imp.(MetricsProvider).Meter(ctx, metric, args...)
}

func (obs *Observer) ReportError(e error) {
	imp, found := obs.Provider.Find(TypeErrorReporter)
	if !found {
		return
	}
	imp.(ErrorReportingProvider).ReportError(e)
}

func Log(msg string, keyValuePairs ...string) {
	observer.Log(msg, keyValuePairs...)
}

func LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
	observer.LogWithLevel(lvl, msg, keyValuePairs...)
}

// Meter logs args to a metrics log from where the values can be aggregated and analyzed.
func Meter(ctx context.Context, metric string, args ...string) {
	observer.Meter(ctx, metric, args...)
}

// ReportError reports error e using the current platform's error reporting provider
func ReportError(e error) {
	observer.ReportError(e)
}
