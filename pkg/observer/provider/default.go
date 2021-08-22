package provider

import (
	"context"
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

	_ GenericProvider = (*defaultProviderImpl)(nil)

	_ ErrorReportingProvider = (*defaultProviderImpl)(nil)
	_ LoggingProvider        = (*defaultProviderImpl)(nil)
	_ MetricsProvider        = (*defaultProviderImpl)(nil)
)

// a NULL provider that does nothing but prevents NPEs in case someone forgets to actually initializa the 'real' provider
func NewDefaultProvider() interface{} {
	return &defaultProviderImpl{}
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
