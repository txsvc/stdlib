package observer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitObserver(t *testing.T) {
	assert.NotNil(t, observer)

	p, found := observer.Provider.Find(TypeLogger)
	assert.True(t, found)
	assert.NotNil(t, p)

	p, found = observer.Provider.Find(TypeErrorReporter)
	assert.True(t, found)
	assert.NotNil(t, p)

	p, found = observer.Provider.Find(TypeMetrics)
	assert.True(t, found)
	assert.NotNil(t, p)
}
