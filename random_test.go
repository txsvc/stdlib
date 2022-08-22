package stdlib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRand(t *testing.T) {
	Seed(time.Now().UnixNano())
	assert.True(t, len(RandString()) > 0)
	assert.True(t, len(RandStringId()) > 0)

	const size = 10
	assert.True(t, len(RandStringN(size)) == size)
}

func BenchmarkRandString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandStringN(10)
	}
}
