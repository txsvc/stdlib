package observer

import (
	"context"

	"github.com/txsvc/stdlib/pkg/provider"
)

const (
	DefaultLogId = "default"
	MetricsLogId = "metric"
	ValuesLogId  = "values"

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
		ReportError(error) error
	}

	MetricsProvider interface {
		Meter(ctx context.Context, metric string, vals ...string)
	}

	LoggingProvider interface {
		Log(string, ...string)
		LogWithLevel(Severity, string, ...string)

		EnableLogging()
		DisableLogging()
	}
)

var (
	p *provider.Provider
)

func NewConfig(opts ...provider.ProviderConfig) (*provider.Provider, error) {
	o, err := provider.New(opts...)
	if err != nil {
		return nil, err
	}
	p = o

	return o, nil
}

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

func EnableLogging() {
	imp, found := p.Find(TypeLogger)
	if !found {
		return
	}
	imp.(LoggingProvider).EnableLogging()
}

func DisableLogging() {
	imp, found := p.Find(TypeLogger)
	if !found {
		return
	}
	imp.(LoggingProvider).DisableLogging()
}

func Meter(ctx context.Context, metric string, vals ...string) {
	imp, found := p.Find(TypeMetrics)
	if !found {
		return
	}
	imp.(MetricsProvider).Meter(ctx, metric, vals...)
}

func ReportError(e error) error {
	imp, found := p.Find(TypeErrorReporter)
	if !found {
		return nil
	}
	return imp.(ErrorReportingProvider).ReportError(e)
}
