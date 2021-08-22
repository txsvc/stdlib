package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	TestProviderImpl struct {
	}
)

func newTestProvider() interface{} {
	return &TestProviderImpl{}
}

func TestWithProvider(t *testing.T) {
	opt := WithProvider("test", TypeLogger, newTestProvider)
	assert.NotNil(t, opt)

	assert.Equal(t, "test", opt.ID)
	assert.Equal(t, TypeLogger, opt.Type)
	assert.NotNil(t, opt.Impl)
}

func TestProviderTypeToString(t *testing.T) {
	assert.Equal(t, "LOGGER", TypeLogger.String())
	assert.Equal(t, "ERROR_REPORTER", TypeErrorReporter.String())
	assert.Equal(t, "METRICS", TypeMetrics.String())
}

func TestNewDefaultProvider(t *testing.T) {
	p := newTestProvider()
	assert.NotNil(t, p)
}
