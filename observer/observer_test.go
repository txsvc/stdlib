package observer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitObserver(t *testing.T) {
	assert.NotNil(t, p)

	pi, found := p.Find(TypeLogger)
	assert.True(t, found)
	assert.NotNil(t, pi)

	pi, found = p.Find(TypeErrorReporter)
	assert.True(t, found)
	assert.NotNil(t, pi)

	pi, found = p.Find(TypeMetrics)
	assert.True(t, found)
	assert.NotNil(t, pi)
}
