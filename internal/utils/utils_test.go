package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestUtils_GetExecuteMode(t *testing.T) {
	t.Run("環境変数がtest指定", func(t *testing.T) {
		os.Setenv("MY_CONTRACTS_EXECUTE_MODE", "test")
		mode, err := GetExecuteMode()
		assert.NoError(t, err)
		assert.Equal(t, Test, mode)
	})
	t.Run("環境変数がdevelopment指定", func(t *testing.T) {
		os.Setenv("MY_CONTRACTS_EXECUTE_MODE", "development")
		mode, err := GetExecuteMode()
		assert.NoError(t, err)
		assert.Equal(t, Development, mode)
	})
	t.Run("環境変数がproduction指定", func(t *testing.T) {
		os.Setenv("MY_CONTRACTS_EXECUTE_MODE", "production")
		mode, err := GetExecuteMode()
		assert.NoError(t, err)
		assert.Equal(t, Production, mode)
	})
	t.Run("環境変数が意味不明指定", func(t *testing.T) {
		os.Setenv("MY_CONTRACTS_EXECUTE_MODE", "aaafdsafdsa")
		mode, err := GetExecuteMode()
		assert.Error(t, err)
		assert.Equal(t, NotDefined, mode)
	})
	t.Run("環境変数が未指定の時", func(t *testing.T) {
		t.Run("go test 実行では testになる", func(t *testing.T) {
			os.Unsetenv("MY_CONTRACTS_EXECUTE_MODE")
			mode, err := GetExecuteMode()
			assert.NoError(t, err)
			assert.Equal(t, Test, mode)
		})
		// go testでは検証できないのでスキップ
		//t.Run("go run 実行ではdevelopmentになる", func(t *testing.T) {
		//	os.Unsetenv("MY_CONTRACTS_EXECUTE_MODE")
		//	mode, err := GetExecuteMode()
		//	assert.NoError(t, err)
		//	assert.Equal(t, Development, mode)
		//})
	})
}
func TestUtils_CreateJstLocation(t *testing.T) {
	expect := time.FixedZone("Asia/Tokyo", 9*60*60)
	assert.EqualValues(t, expect, CreateJstLocation())
}
