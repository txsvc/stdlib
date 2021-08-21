package loader

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func LoaderTestFunc(ctx context.Context, key string) (interface{}, error) {
	return "OK", nil
}

func TestNew(t *testing.T) {
	ld, err := New(nil, DefaultTTL)
	assert.Error(t, err, MsgMissingLoaderFunc)
	assert.Nil(t, ld)

	ld, err = New(LoaderTestFunc, -1)
	assert.Error(t, err, MsgInvalidTTL)
	assert.Nil(t, ld)

	ld, err = New(LoaderTestFunc, 0)
	assert.Nil(t, err)
	assert.NotNil(t, ld)

	ld, err = New(LoaderTestFunc, DefaultTTL)
	assert.Nil(t, err)
	assert.NotNil(t, ld)
}

func TestSimpleLoad(t *testing.T) {
	ld, _ := New(LoaderTestFunc, DefaultTTL)

	item1, err := ld.Load(context.TODO(), "key1")
	assert.Nil(t, err)
	assert.NotNil(t, item1)

	// test contains
	assert.True(t, ld.Contains(context.TODO(), "key1"))
	assert.False(t, ld.Contains(context.TODO(), "other_key"))

	// test remove
	ld.Remove(context.TODO(), "key1")
	assert.False(t, ld.Contains(context.TODO(), "key1"))
}

func TestUpdate(t *testing.T) {
	ld, _ := New(LoaderTestFunc, DefaultTTL)

	item1, err := ld.Load(context.TODO(), "key1")
	assert.Nil(t, err)
	assert.NotNil(t, item1)
	assert.Equal(t, "OK", item1)

	// update
	err = ld.Update(context.TODO(), "key1", "NEW")
	assert.Nil(t, err)

	// verify
	item1, err = ld.Load(context.TODO(), "key1")
	assert.Nil(t, err)
	assert.NotNil(t, item1)
	assert.Equal(t, "NEW", item1)
}

func TestStats(t *testing.T) {
	ld, _ := New(LoaderTestFunc, DefaultTTL)

	item1, err := ld.Load(context.TODO(), "key1")
	assert.Nil(t, err)
	assert.NotNil(t, item1)

	// assume a cache miss at the first load
	assert.Equal(t, int64(0), ld.Hits())
	assert.Equal(t, int64(1), ld.Misses())

	// assume a cache hit the next time
	item1, err = ld.Load(context.TODO(), "key1")
	assert.Nil(t, err)
	assert.NotNil(t, item1)
	assert.Equal(t, int64(1), ld.Hits())
	assert.Equal(t, int64(1), ld.Misses())

	// ratio 50% after 1 miss, 1 hit
	assert.Equal(t, float64(0.5), ld.Ratio())
	// ratio 66% after 1 miss, 2 hits
	ld.Load(context.TODO(), "key1")
	assert.Greater(t, ld.Ratio(), float64(0.6))

	// stats string
	stats := ld.Stats()
	assert.NotEmpty(t, stats)

	parts := strings.Split(stats, ",")
	assert.Equal(t, 2, len(parts))
	assert.Equal(t, "2", parts[0]) // 2 hits
	assert.Equal(t, "1", parts[1]) // 1 miss
}
