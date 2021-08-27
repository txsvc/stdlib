package observer

import (
	"context"
	"log"

	"github.com/txsvc/stdlib/pkg/provider"
)

type (
	// defaultObserverImpl provides a simple implementation in the absence of any configuration
	defaultObserverImpl struct {
	}
)

var (
	// Interface guards.

	// This enforces a compile-time check of the provider implmentation,
	// making sure all the methods defined in the interfaces are implemented.

	_ provider.GenericProvider = (*defaultObserverImpl)(nil)

	_ ErrorReportingProvider = (*defaultObserverImpl)(nil)
	_ LoggingProvider        = (*defaultObserverImpl)(nil)
	_ MetricsProvider        = (*defaultObserverImpl)(nil)

	// the instance, a singleton
	theDefaultProvider *defaultObserverImpl
)

func init() {
	theDefaultProvider = &defaultObserverImpl{}
	reset()
}

func reset() {
	// initialize the observer with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real provider
	loggingConfig := provider.WithProvider("observer.default.logger", TypeLogger, NewDefaultProvider)
	errorReportingConfig := provider.WithProvider("observer.default.errorreporting", TypeErrorReporter, NewDefaultProvider)
	metricsConfig := provider.WithProvider("observer.default.metrics", TypeMetrics, NewDefaultProvider)

	o, err := provider.New()
	if err != nil {
		log.Fatal(err)
	}

	err = o.RegisterProviders(false, loggingConfig, errorReportingConfig, metricsConfig)
	if err != nil {
		log.Fatal(err)
	}
	p = o
}

// a default provider that does nothing but prevents NPEs in case someone forgets to actually initializa the 'real' provider
func NewDefaultProvider() interface{} {
	return theDefaultProvider
}

// IF GenericProvider

func (np *defaultObserverImpl) Close() error {
	return nil
}

// IF ErrorReportingProvider

func (np *defaultObserverImpl) ReportError(e error) {
}

// IF LoggingProvider

func (np *defaultObserverImpl) Log(msg string, keyValuePairs ...string) {
}

func (np *defaultObserverImpl) LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
}

// IF MetricsProvider

func (np *defaultObserverImpl) Meter(ctx context.Context, metric string, args ...string) {
}
