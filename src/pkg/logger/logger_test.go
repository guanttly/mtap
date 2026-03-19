package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	l, err := NewLogger("info")
	assert.NoError(t, err)
	assert.NotNil(t, l)
}

func TestNewLogger_InvalidLevel(t *testing.T) {
	l, err := NewLogger("invalid_level")
	assert.NoError(t, err)
	assert.NotNil(t, l) // fallback to info level
}

func TestL_Default(t *testing.T) {
	l := L()
	assert.NotNil(t, l)
}

func TestSetDefault(t *testing.T) {
	original := L()
	newLogger, _ := NewLogger("debug")
	SetDefault(newLogger)
	assert.Equal(t, newLogger, L())
	SetDefault(original) // restore
}

func TestInfoWarnError_NoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		Info("test info", "key", "value")
		Warn("test warn", "key", "value")
		Error("test error", "key", "value")
	})
}
