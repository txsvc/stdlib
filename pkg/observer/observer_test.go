package observer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/txsvc/stdlib/pkg/observer/provider"
)

type (
	TestProviderImpl struct {
	}
)

func newTestProvider() interface{} {
	return &TestProviderImpl{}
}

func (np *TestProviderImpl) ReportError(e error) {
}

func (np *TestProviderImpl) Log(msg string, keyValuePairs ...string) {
}

func (np *TestProviderImpl) LogWithLevel(lvl provider.Severity, msg string, keyValuePairs ...string) {
}

func (np *TestProviderImpl) Meter(ctx context.Context, metric string, args ...string) {
}

func TestWithProvider(t *testing.T) {
	opt := provider.WithProvider("test", provider.TypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	assert.Equal(t, "test", opt.ID)
	assert.Equal(t, provider.TypeLogger, opt.Type)
	assert.NotNil(t, opt.Impl)
}

func TestInitDefaultObserver(t *testing.T) {
	reset()

	p := DefaultObserver()

	assert.NotNil(t, p)
	assert.NotNil(t, p.errorReportingProvider)
	assert.NotNil(t, p.loggingProvider)
	assert.NotNil(t, p.metricsProvdider)
}

func TestInitObserver(t *testing.T) {
	reset()

	p, err := Init(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, p)

	assert.Equal(t, 0, len(p.providers))

	assert.Nil(t, p.errorReportingProvider)
	assert.Nil(t, p.loggingProvider)
	assert.Nil(t, p.metricsProvdider)
}

func TestInitObserverDuplicateProvider(t *testing.T) {
	reset()

	opt1 := provider.WithProvider("test1", provider.TypeLogger, newTestProvider)
	opt2 := provider.WithProvider("test2", provider.TypeLogger, newTestProvider)

	p, err := Init(context.Background(), opt1, opt2)
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestRegisterObserver(t *testing.T) {
	reset()

	defaultPlatform := DefaultObserver()
	p, err := Init(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, defaultPlatform)
	assert.NotNil(t, p)

	assert.NotEqual(t, defaultPlatform, p)

	old := Register(p)
	assert.NotNil(t, old)
	assert.Equal(t, defaultPlatform, old)

	newDefaultPlatform := DefaultObserver()
	assert.NotNil(t, newDefaultPlatform)
	assert.Equal(t, newDefaultPlatform, p)
}

func TestRegisterProvider(t *testing.T) {
	reset()

	opt := provider.WithProvider("test", provider.TypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	p := DefaultObserver()
	assert.NotNil(t, p)

	err := p.RegisterProviders(false, opt)
	assert.Error(t, err)

	err = p.RegisterProviders(true, opt)
	assert.NoError(t, err)

}

func TestGetDefaultProviders(t *testing.T) {
	reset()

	p1, ok := Provider(provider.TypeLogger)
	assert.True(t, ok)
	assert.NotNil(t, p1)

	logger := p1.(provider.LoggingProvider)
	assert.NotNil(t, logger)

	p2, ok := Provider(provider.TypeErrorReporter)
	assert.True(t, ok)
	assert.NotNil(t, p2)

	errorReporter := p2.(provider.ErrorReportingProvider)
	assert.NotNil(t, errorReporter)

	p3, ok := Provider(provider.TypeMetrics)
	assert.True(t, ok)
	assert.NotNil(t, p3)

	metrics := p3.(provider.MetricsProvider)
	assert.NotNil(t, metrics)
}
