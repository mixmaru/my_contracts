package my_logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestMyLogger_GetLogger(t *testing.T) {
	t.Run("シングルトン的に1つのloggerが取得されることを確認", func(t *testing.T) {
		logger1, err := GetLogger()
		assert.NoError(t, err)
		logger2, err := GetLogger()
		assert.NoError(t, err)

		assert.Equal(t, logger1, logger2)
	})

	t.Run("production modeのとき", func(t *testing.T) {
		// 本番モードを再現
		err := os.Setenv("MY_CONTRACTS_EXECUTE_MODE", "production")
		assert.NoError(t, err)

		logger, err := GetLogger()
		assert.NoError(t, err)
		assert.IsType(t, &zap.Logger{}, logger)
	})

	t.Run("development modeのとき", func(t *testing.T) {
		// developmentモードを再現
		err := os.Setenv("MY_CONTRACTS_EXECUTE_MODE", "development")
		assert.NoError(t, err)

		logger, err := GetLogger()
		assert.NoError(t, err)
		assert.IsType(t, &zap.Logger{}, logger)
	})

	t.Run("test modeのとき", func(t *testing.T) {
		// testモードを再現
		err := os.Setenv("MY_CONTRACTS_EXECUTE_MODE", "test")
		assert.NoError(t, err)

		logger, err := GetLogger()
		assert.NoError(t, err)
		assert.IsType(t, &zap.Logger{}, logger)
	})
}
