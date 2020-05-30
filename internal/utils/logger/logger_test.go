package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger_GetLogger(t *testing.T) {
	logger1, err := GetLogger()
	assert.NoError(t, err)
	logger2, err := GetLogger()
	assert.NoError(t, err)

	assert.Equal(t, logger1, logger2)
}
