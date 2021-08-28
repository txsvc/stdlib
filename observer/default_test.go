package observer

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	assert.NotNil(t, p)
	Log("info message")
}

func TestWithLogLevel(t *testing.T) {
	assert.NotNil(t, p)
	LogWithLevel(LevelWarn, "LevelWarn message")
	LogWithLevel(LevelInfo, "LevelInfo message")
	LogWithLevel(LevelError, "LevelError message")
	LogWithLevel(LevelDebug, "LevelDebug message")
}

func TestOffAndOn(t *testing.T) {
	assert.NotNil(t, p)

	DisableLogging()
	Log("should NOT see this message")

	EnableLogging()
	Log("SHOULD see this message")
}

func TestReportError(t *testing.T) {
	assert.NotNil(t, p)

	// simple exception
	e := fmt.Errorf("an error happened")
	ReportError(e)
	// with stacktrace
	ReportError(outer())
}

func inner() error {
	return errors.New("seems we have an error here")
}

func middle() error {
	err := inner()
	if err != nil {
		return err
	}
	return nil
}

func outer() error {
	err := middle()
	if err != nil {
		return err
	}
	return nil
}
