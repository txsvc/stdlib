package stdlib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMnemonicGenerateNew(t *testing.T) {
	// Test generating a new mnemonic with empty phrase
	mnemonic, err := CreateMnemonic("")
	assert.NoError(t, err)
	assert.NotEmpty(t, mnemonic)

	// Check that we have at least MinWordsInPassPhrase words
	wordCount := strings.Count(mnemonic, " ")
	assert.GreaterOrEqual(t, wordCount, MinWordsInPassPhrase)

	// Verify it's not empty and contains only valid characters
	assert.True(t, len(mnemonic) > 0)
	assert.False(t, strings.Contains(mnemonic, "  ")) // No double spaces
	assert.False(t, strings.HasPrefix(mnemonic, " ")) // No leading space
	assert.False(t, strings.HasSuffix(mnemonic, " ")) // No trailing space
}

func TestCreateMnemonicValidateExistingPhrase(t *testing.T) {
	// Test with a valid phrase that has enough words
	validPhrase := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	mnemonic, err := CreateMnemonic(validPhrase)
	assert.NoError(t, err)
	assert.Equal(t, validPhrase, mnemonic)
}

func TestCreateMnemonicNormalizeWithExtraSpaces(t *testing.T) {
	// Test phrase normalization - extra spaces should be removed
	phraseWithSpaces := "  abandon   abandon  abandon   abandon  abandon  abandon abandon abandon abandon abandon abandon about  "
	expectedPhrase := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	mnemonic, err := CreateMnemonic(phraseWithSpaces)
	assert.NoError(t, err)
	assert.Equal(t, expectedPhrase, mnemonic)
}

func TestCreateMnemonicNormalizeWithEmptyWords(t *testing.T) {
	// Test phrase normalization - empty words should be filtered out
	phraseWithEmptyWords := "abandon  abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	expectedPhrase := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	mnemonic, err := CreateMnemonic(phraseWithEmptyWords)
	assert.NoError(t, err)
	assert.Equal(t, expectedPhrase, mnemonic)
}

func TestCreateMnemonicRejectPhraseTooShort(t *testing.T) {
	// Test with a phrase that's too short (less than MinWordsInPassPhrase)
	shortPhrase := "abandon abandon abandon"

	mnemonic, err := CreateMnemonic(shortPhrase)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassPhrase, err)
	assert.Empty(t, mnemonic)
}

func TestCreateMnemonicRejectEmptyPhraseAfterNormalization(t *testing.T) {
	// Test with a phrase that becomes too short after normalization
	phraseWithOnlySpaces := "   "

	mnemonic, err := CreateMnemonic(phraseWithOnlySpaces)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidPassPhrase, err)
	assert.Empty(t, mnemonic)
}

func TestCreateMnemonicMinimumWordsBoundary(t *testing.T) {
	// Test exactly at the boundary (MinWordsInPassPhrase words)
	words := make([]string, MinWordsInPassPhrase+1) // +1 because count is spaces, not words
	for i := range words {
		words[i] = "word"
	}
	boundaryPhrase := strings.Join(words, " ")

	mnemonic, err := CreateMnemonic(boundaryPhrase)
	assert.NoError(t, err)
	assert.Equal(t, boundaryPhrase, mnemonic)
}

func TestCreateMnemonicGeneratedMnemonicsAreDifferent(t *testing.T) {
	// Test that generated mnemonics are different (randomness)
	mnemonic1, err1 := CreateMnemonic("")
	assert.NoError(t, err1)

	mnemonic2, err2 := CreateMnemonic("")
	assert.NoError(t, err2)

	// They should be different (extremely unlikely to be the same)
	assert.NotEqual(t, mnemonic1, mnemonic2)
}
