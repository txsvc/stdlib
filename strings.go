package stdlib

// FIXME: move to stdlib !

// This code is adopted from stackoverflow.com. Use with caution.
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go

import (
	crand "crypto/rand"
	"fmt"
	"strings"
)

const (
	// Used in RandString
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	// Same as above only for RandStringSimple
	simpleBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

	// Used to generate passwords. this 'sacrifies' some letters in order to make room for other characters.
	passwordBytes = "abcdefghijklmnpqrsuvwyzABCDEFGIJKLNOPQRSTUVWXYZ0123456789#%*+@_-"

	idLen          = 8
	defaultRandLen = 8
)

// RandString returns a random string.
func RandString() string {
	return RandStringN(defaultRandLen)
}

// RandString returns a random sting of lenght n based on a..zA..z0..9
func RandStringN(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
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

// RandStringSimple returns a random sting of lenght n based on a..z0..9.
// This is basically the same as RandStringN but only case-insensitive.
func RandStringSimple(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(simpleBytes) {
			sb.WriteByte(simpleBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}

func RandPasswordString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(passwordBytes) {
			sb.WriteByte(passwordBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
