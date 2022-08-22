package loader

// Dataloader inspired by https://github.com/facebook/dataloader

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	mcache "github.com/OrlovEvgeny/go-mcache"
)

const (
	// DefaultTTL is the default TTL used if nothing else is specified.
	// Keep this short by default to support subsequent queries of the same item
	// and to avoid a stale cache in the long run.
	DefaultTTL = time.Minute * 2

	MsgInvalidTTL        = "TTL can not be less than zero"
	MsgMissingLoaderFunc = "loader function is missing"
)

var (
	ErrInvalidTTL        = errors.New(MsgInvalidTTL)
	ErrMissingLoaderFunc = errors.New(MsgMissingLoaderFunc)
)

type (
	// LoaderFunc abstracts the process of loading a resource
	LoaderFunc func(context.Context, string) (interface{}, error)

	// Loader holds cached resources. The cache is a simple in-memory cache with TTL.
	Loader struct {
		load         LoaderFunc
		cache        *mcache.CacheDriver
		expiresAfter time.Duration
		mu           sync.Mutex
		cacheHit     int64
		cacheMiss    int64
	}
)

// New initializes the loader
func New(lf LoaderFunc, ttl time.Duration) (*Loader, error) {
	if ttl < 0 {
		return nil, ErrInvalidTTL
	}
	if lf == nil {
		return nil, ErrMissingLoaderFunc
	}
	return &Loader{
		load:         lf,
		cache:        mcache.New(),
		expiresAfter: ttl,
	}, nil
}

// Load returns either a cached resource or calls the loader function to retrieve the requested resource
func (ld *Loader) Load(ctx context.Context, key string) (interface{}, error) {

	ld.mu.Lock()
	defer ld.mu.Unlock()

	if data, ok := ld.cache.Get(key); ok {
		ld.cacheHit++
		return data, nil
	}

	data, err := ld.load(ctx, key)
	if err != nil {
		return nil, err
	}
	if data != nil {
		if err := ld.cache.Set(key, data, ld.expiresAfter); err != nil {
			return nil, err
		}
		ld.cacheMiss++
		return data, nil
	}
	ld.cacheMiss++
	return nil, nil
}

func (ld *Loader) Contains(ctx context.Context, key string) bool {
	_, ok := ld.cache.Get(key)
	return ok
}

// Update replaces a resource in the cache with a different one without loading it. It also resets the TTL.
func (ld *Loader) Update(ctx context.Context, key string, value interface{}) error {
	ld.mu.Lock()
	defer ld.mu.Unlock()

	if err := ld.cache.Set(key, value, ld.expiresAfter); err != nil {
		return err
	}
	return nil
}

// Remove removes a resource from the cache if it is there. The function does nothing otherwise.
func (ld *Loader) Remove(ctx context.Context, key string) {
	ld.mu.Lock()
	ld.cache.Remove(key)
	ld.mu.Unlock()
}

// some metrics

func (ld *Loader) Hits() int64 {
	return ld.cacheHit
}

func (ld *Loader) Misses() int64 {
	return ld.cacheMiss
}

func (ld *Loader) Ratio() float64 {
	total := ld.cacheHit + ld.cacheMiss
	if total == 0 {
		return 0
	}
	return float64(ld.cacheHit) / float64(total)
}

func (ld *Loader) Stats() string {
	return fmt.Sprintf("%d,%d", ld.cacheHit, ld.cacheMiss)
}
