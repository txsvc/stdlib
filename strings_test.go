package stdlib

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	s := RandStringN(10)
	assert.NotEmpty(t, s)
	assert.Equal(t, 10, len(s))

	s = RandStringN(1024)
	assert.NotEmpty(t, s)
	assert.Equal(t, 1024, len(s))

	// just make this long enough that it is not likely that anything is missing. could randomly fail though.

	containsAllChars := true
	for i := 0; i < len(letterBytes); i++ {
		if !strings.Contains(s, string(letterBytes[i])) {
			containsAllChars = false
			fmt.Printf("%s: missing: '%s'\n", s, string(letterBytes[i]))
			break
		}
	}
	assert.True(t, containsAllChars)
}

func TestRandStringSimple(t *testing.T) {
	s := RandStringSimple(16)
	assert.NotEmpty(t, s)
	assert.Equal(t, 16, len(s))

	s = RandStringSimple(64)
	assert.NotEmpty(t, s)
	assert.Equal(t, 64, len(s))

	fmt.Println(s)
}

func TestRandPasswordString(t *testing.T) {
	s := RandPasswordString(16)
	assert.NotEmpty(t, s)
	assert.Equal(t, 16, len(s))

	s = RandPasswordString(64)
	assert.NotEmpty(t, s)
	assert.Equal(t, 64, len(s))

	fmt.Println(s)
}

func TestRandStringZeroLength(t *testing.T) {
	s := RandStringN(0)
	assert.Empty(t, s)
	assert.Equal(t, 0, len(s))
}

func TestRandStringSingleCharacter(t *testing.T) {
	s := RandStringN(1)
	assert.NotEmpty(t, s)
	assert.Equal(t, 1, len(s))
	// Should be one of the valid characters
	assert.Contains(t, letterBytes, s)
}

func TestRandStringLargeString(t *testing.T) {
	s := RandStringN(10000)
	assert.NotEmpty(t, s)
	assert.Equal(t, 10000, len(s))
}

func TestRandStringSimpleZeroLength(t *testing.T) {
	s := RandStringSimple(0)
	assert.Empty(t, s)
	assert.Equal(t, 0, len(s))
}

func TestRandStringSimpleSingleCharacter(t *testing.T) {
	s := RandStringSimple(1)
	assert.NotEmpty(t, s)
	assert.Equal(t, 1, len(s))
	// Should be one of the valid simple characters
	assert.Contains(t, simpleBytes, s)
}

func TestRandStringSimpleVerifyCharacterSet(t *testing.T) {
	s := RandStringSimple(1000)

	// Verify only lowercase letters and digits are used
	for _, char := range s {
		assert.True(t, (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9'),
			"Character '%c' is not in the expected simple character set", char)
	}
}

func TestRandPasswordStringZeroLength(t *testing.T) {
	s := RandPasswordString(0)
	assert.Empty(t, s)
	assert.Equal(t, 0, len(s))
}

func TestRandPasswordStringSingleCharacter(t *testing.T) {
	s := RandPasswordString(1)
	assert.NotEmpty(t, s)
	assert.Equal(t, 1, len(s))
	// Should be one of the valid password characters
	assert.Contains(t, passwordBytes, s)
}

func TestRandPasswordStringVerifySpecialCharacters(t *testing.T) {
	// Generate a long password and verify it contains special characters
	s := RandPasswordString(1000)

	// Should contain at least some special characters from the password set
	hasSpecial := false
	specialChars := "#%*+@_-"
	for _, char := range s {
		if strings.Contains(specialChars, string(char)) {
			hasSpecial = true
			break
		}
	}
	// With 1000 characters, we should statistically have special chars
	assert.True(t, hasSpecial, "Password should contain special characters")
}

func TestLetterBytesContainsExpected(t *testing.T) {
	// Verify letterBytes contains expected characters
	assert.Contains(t, letterBytes, "a")
	assert.Contains(t, letterBytes, "z")
	assert.Contains(t, letterBytes, "A")
	assert.Contains(t, letterBytes, "Z")
	assert.Contains(t, letterBytes, "0")
	assert.Contains(t, letterBytes, "9")
}

func TestSimpleBytesIsSubset(t *testing.T) {
	// Verify simpleBytes is a subset of letterBytes (lowercase + digits only)
	for _, char := range simpleBytes {
		assert.True(t, strings.Contains(letterBytes, string(char)),
			"Simple character '%c' should be in letterBytes", char)
	}
}

func TestPasswordBytesContainsSpecial(t *testing.T) {
	// Verify passwordBytes contains special characters
	assert.Contains(t, passwordBytes, "#")
	assert.Contains(t, passwordBytes, "%")
	assert.Contains(t, passwordBytes, "*")
	assert.Contains(t, passwordBytes, "+")
	assert.Contains(t, passwordBytes, "@")
	assert.Contains(t, passwordBytes, "_")
	assert.Contains(t, passwordBytes, "-")
}
