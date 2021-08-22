package observer

import (
	"context"
	"fmt"
	"log"

	"github.com/txsvc/stdlib/pkg/observer/provider"
)

const (
	MsgMissingProvider = "provider '%s' required"
)

type (
	Observer struct {
		errorReportingProvider provider.ErrorReportingProvider
		metricsProvdider       provider.MetricsProvider
		loggingProvider        provider.LoggingProvider

		providers map[provider.ProviderType]provider.ProviderConfig
	}
)

var (
	// internal
	observer *Observer
)

func init() {
	reset()
}

func reset() {
	// initialize the observer with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real provider
	loggingConfig := provider.WithProvider("observer.null.logger", provider.TypeLogger, provider.NewDefaultProvider)
	errorReportingConfig := provider.WithProvider("observer.null.errorreporting", provider.TypeErrorReporter, provider.NewDefaultProvider)
	metricsConfig := provider.WithProvider("observer.null.metrics", provider.TypeMetrics, provider.NewDefaultProvider)

	p, err := Init(context.Background(), loggingConfig, errorReportingConfig, metricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	Register(p)
}

// Init creates a new global observer instance and configures it with providers
func Init(ctx context.Context, opts ...provider.ProviderConfig) (*Observer, error) {
	p := Observer{
		providers: make(map[provider.ProviderType]provider.ProviderConfig),
	}

	if err := p.RegisterProviders(false, opts...); err != nil {
		return nil, err
	}

	return &p, nil
}

// Register makes o the new default observer
func Register(o *Observer) *Observer {
	if o == nil {
		return nil
	}
	old := observer
	observer = o
	return old
}

// RegisterProviders registers one or more providers. An existing provider will be overwritten
// if ignoreExists is true, otherwise the function returns an error.
func (p *Observer) RegisterProviders(ignoreExists bool, opts ...provider.ProviderConfig) error {
	for _, opt := range opts {
		if _, ok := p.providers[opt.Type]; ok {
			if !ignoreExists {
				return fmt.Errorf("provider of type '%s' already registered", opt.Type.String())
			}
		}
		p.providers[opt.Type] = opt

		switch opt.Type {
		case provider.TypeLogger:
			p.loggingProvider = opt.Impl().(provider.LoggingProvider)
		case provider.TypeErrorReporter:
			p.errorReportingProvider = opt.Impl().(provider.ErrorReportingProvider)
		case provider.TypeMetrics:
			p.metricsProvdider = opt.Impl().(provider.MetricsProvider)
		}
	}
	return nil
}

// DefaultObserver returns the current default observer.
func DefaultObserver() *Observer {
	return observer
}

// Provider returns the registered provider instance if it is defined.
// The bool flag is set to true if there is a provider and false otherwise.
func Provider(providerType provider.ProviderType) (interface{}, bool) {
	opt, ok := observer.providers[providerType]
	if !ok {
		return nil, false
	}
	return opt.Impl(), true
}

// a set of convenience functions in order to avoid getting the provider impl every time

func Log(msg string, keyValuePairs ...string) {
	observer.loggingProvider.Log(msg, keyValuePairs...)
}

func LogWithLevel(lvl provider.Severity, msg string, keyValuePairs ...string) {
	observer.loggingProvider.LogWithLevel(lvl, msg, keyValuePairs...)
}

// Meter logs args to a metrics log from where the values can be aggregated and analyzed.
func Meter(ctx context.Context, metric string, args ...string) {
	observer.metricsProvdider.Meter(ctx, metric, args...)
}

// ReportError reports error e using the current platform's error reporting provider
func ReportError(e error) {
	observer.errorReportingProvider.ReportError(e)
}
