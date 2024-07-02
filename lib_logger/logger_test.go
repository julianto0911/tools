package lib_logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	//set logger with debug false
	l, err := NewLogger("C:/work/logs/", false)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, l, "object should be initiated")

	//set logger with debug true
	l, err = NewLogger("C:/work/logs/", true)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, l, "object should be initiated")
}

func TestMockLogs(t *testing.T) {
	t.Run("Handle log message with fields", func(t *testing.T) {
		// Given
		logs, obs := MockLogs()
		// When
		myFunction(logs)

		// Then
		require.Equal(t, 2, obs.Len())
		allLogs := obs.All()
		assert.Equal(t, "log myFunction", allLogs[0].Message)
		assert.Equal(t, "log with fields", allLogs[1].Message)
		assert.ElementsMatch(t, []zap.Field{
			{Key: "keyOne", String: "valueOne"},
			{Key: "keyTwo", String: "valueTwo"},
		}, allLogs[1].Context)
	})
}

func myFunction(logger *zap.Logger) {
	logger.Info("log myFunction")
	logger.With(
		zap.Field{Key: "keyOne", String: "valueOne"},
		zap.Field{Key: "keyTwo", String: "valueTwo"},
	).Info("log with fields")
}
