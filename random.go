package stdlib

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits  = 6 // 6 bits to represent a letter index
	idLen          = 8
	defaultRandLen = 8
	letterIdxMask  = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax   = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = newLockedSource(time.Now().UnixNano())

type lockedSource struct {
	source rand.Source
	lock   sync.Mutex
}

// newLockedSource creates a new thread-safe random source with the given seed.
// It wraps a rand.Source with a mutex to ensure concurrent access is safe.
func newLockedSource(seed int64) *lockedSource {
	return &lockedSource{
		source: rand.NewSource(seed),
	}
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64
// in a thread-safe manner.
func (ls *lockedSource) Int63() int64 {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	return ls.source.Int63()
}

// Seed uses the provided seed value to initialize the generator to a deterministic state
// in a thread-safe manner.
func (ls *lockedSource) Seed(seed int64) {
	ls.lock.Lock()
	defer ls.lock.Unlock()
	ls.source.Seed(seed)
}

// RandString returns a random string.
func RandString() string {
	return RandStringN(defaultRandLen)
}

// RandStringId returns a random id string.
func RandStringId() string {
	b := make([]byte, idLen)
	_, err := crand.Read(b)
	if err != nil {
		return RandStringN(idLen)
	}

	return fmt.Sprintf("%x%x%x%x", b[0:2], b[2:4], b[4:6], b[6:8])
}

// Randn returns a random string with length n.
func RandStringN(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// Seed sets the seed for the default random source.
// This affects all random string generation functions that don't use crypto/rand.
func Seed(seed int64) {
	src.Seed(seed)
}
