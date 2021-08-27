package observer

import (
	"context"
	"log"

	"github.com/txsvc/stdlib/pkg/provider"
)

type (
	// defaultProviderImpl provides a simple implementation in the absence of any configuration
	defaultProviderImpl struct {
	}
)

var (
	// Interface guards.

	// This enforces a compile-time check of the provider implmentation,
	// making sure all the methods defined in the interfaces are implemented.

	_ provider.GenericProvider = (*defaultProviderImpl)(nil)

	_ ErrorReportingProvider = (*defaultProviderImpl)(nil)
	_ LoggingProvider        = (*defaultProviderImpl)(nil)
	_ MetricsProvider        = (*defaultProviderImpl)(nil)

	// the instance, a singleton
	theDefaultProvider *defaultProviderImpl
)

func init() {
	theDefaultProvider = &defaultProviderImpl{}
	reset()
}

func reset() {
	// initialize the observer with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real provider
	loggingConfig := provider.WithProvider("observer.null.logger", TypeLogger, NewDefaultProvider)
	errorReportingConfig := provider.WithProvider("observer.null.errorreporting", TypeErrorReporter, NewDefaultProvider)
	metricsConfig := provider.WithProvider("observer.null.metrics", TypeMetrics, NewDefaultProvider)

	p, err := provider.New()
	if err != nil {
		log.Fatal(err)
	}
	o := &Observer{
		Provider: p,
	}

	err = o.Provider.RegisterProviders(false, loggingConfig, errorReportingConfig, metricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	observer = o
}

// a NULL provider that does nothing but prevents NPEs in case someone forgets to actually initializa the 'real' provider
func NewDefaultProvider() interface{} {
	return theDefaultProvider
}

// IF GenericProvider

func (np *defaultProviderImpl) Close() error {
	return nil
}

// IF ErrorReportingProvider

func (np *defaultProviderImpl) ReportError(e error) {
}

// IF LoggingProvider

func (np *defaultProviderImpl) Log(msg string, keyValuePairs ...string) {
}

func (np *defaultProviderImpl) LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
}

// IF MetricsProvider

func (np *defaultProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
}
