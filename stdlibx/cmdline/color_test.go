package cmdline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithColor(t *testing.T) {
	output := WithColor("Hello", BgRed)
	assert.NotEmpty(t, output)
	//assert.Equal(t, "Hello", output) // FIXME test that we actually have RED as background
}

func TestWithColorPadding(t *testing.T) {
	output := WithColorPadding("Hello", BgRed)
	assert.NotEmpty(t, output)
	//assert.Equal(t, " Hello ", output) // FIXME test that we actually have RED as background
}
