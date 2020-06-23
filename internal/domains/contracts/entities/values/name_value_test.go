package values

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNameValue_NewNameValue(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		name, err := NewNameValue("個人顧客名")
		assert.NoError(t, err)
		assert.Equal(t, NameValue{"個人顧客名"}, name)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		name, err := NewNameValue("")
		assert.Error(t, err)
		assert.Equal(t, NameValue{}, name)
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		name, err := NewNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Error(t, err)
		assert.Equal(t, NameValue{}, name)
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		name, err := NewNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Equal(t, NameValue{"0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789"}, name)
	})
}

func TestNameValue_NameValidate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		errs := NameValidate("個人顧客名")
		assert.Len(t, errs, 0)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		errs := NameValidate("")
		assert.Len(t, errs, 1)
		assert.EqualError(t, errs[0], "空です")
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		errs := NameValidate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Len(t, errs, 1)
		assert.EqualError(t, errs[0], "50文字より多いです")
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		errs := NameValidate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.Len(t, errs, 0)
	})

}
