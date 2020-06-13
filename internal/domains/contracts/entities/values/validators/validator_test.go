package validators

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidator_IsEmptyString(t *testing.T) {
	t.Run("から文字ではない。", func(t *testing.T) {
		assert.False(t, IsEmptyString("a"))
		assert.False(t, IsEmptyString("あ"))
		assert.False(t, IsEmptyString("0"))
		assert.False(t, IsEmptyString("-"))
	})

	t.Run("から文字である。", func(t *testing.T) {
		assert.True(t, IsEmptyString(""))
	})
}

func TestValidator_IsOverLengthString(t *testing.T) {
	t.Run("指定文字数を超えていない", func(t *testing.T) {
		assert.False(t, IsOverLengthString("a", 10))
		assert.False(t, IsOverLengthString("あ", 10))
		assert.False(t, IsOverLengthString("0", 10))
		assert.False(t, IsOverLengthString("-", 10))
	})

	t.Run("指定文字数を超えている", func(t *testing.T) {
		assert.True(t, IsOverLengthString("aaaaaaaaaaa", 10))
		assert.True(t, IsOverLengthString("あああああああああああ", 10))
		assert.True(t, IsOverLengthString("00000000000", 10))
		assert.True(t, IsOverLengthString("-----------", 10))
	})

	t.Run("指定文字数と同じ", func(t *testing.T) {
		assert.False(t, IsOverLengthString("aaaaaaaaaa", 10))
		assert.False(t, IsOverLengthString("ああああああああああ", 10))
		assert.False(t, IsOverLengthString("0000000000", 10))
		assert.False(t, IsOverLengthString("----------", 10))
	})

	t.Run("空文字だったとき", func(t *testing.T) {
		assert.False(t, IsOverLengthString("", 10))
	})
}
