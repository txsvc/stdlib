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

	TypeLogger        provider.ProviderType = 10
	TypeErrorReporter provider.ProviderType = 11
	TypeMetrics       provider.ProviderType = 12
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

var (
	p *provider.Provider
)

func Log(msg string, keyValuePairs ...string) {
	imp, found := p.Find(TypeLogger)
	if !found {
		return
	}
	imp.(LoggingProvider).Log(msg, keyValuePairs...)
}

func LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
	imp, found := p.Find(TypeLogger)
	if !found {
		return
	}
	imp.(LoggingProvider).LogWithLevel(lvl, msg, keyValuePairs...)
}

func Meter(ctx context.Context, metric string, args ...string) {
	imp, found := p.Find(TypeMetrics)
	if !found {
		return
	}
	imp.(MetricsProvider).Meter(ctx, metric, args...)
}

func ReportError(e error) {
	imp, found := p.Find(TypeErrorReporter)
	if !found {
		return
	}
	imp.(ErrorReportingProvider).ReportError(e)
}
