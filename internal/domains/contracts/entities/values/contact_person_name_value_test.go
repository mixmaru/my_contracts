package values

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContactPersonNameValue_NewContactPersonNameValue(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		name, err := NewContactPersonNameValue("担当者名")
		assert.NoError(t, err)
		assert.Equal(t, ContactPersonNameValue{"担当者名"}, name)
	})

	t.Run("担当者名が空文字だった時", func(t *testing.T) {
		name, err := NewContactPersonNameValue("")
		assert.Error(t, err)
		assert.Equal(t, ContactPersonNameValue{}, name)
	})

	t.Run("担当者名が50文字を超えていた時", func(t *testing.T) {
		name, err := NewContactPersonNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Error(t, err)
		assert.Equal(t, ContactPersonNameValue{}, name)
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		name, err := NewContactPersonNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Equal(t, ContactPersonNameValue{"0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789"}, name)
	})
}

func TestNameValue_ContactPersonNameValidate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		errs := ContactPersonNameValidate("担当者名")
		assert.Len(t, errs, 0)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		errs := ContactPersonNameValidate("")
		assert.Len(t, errs, 1)
		assert.EqualError(t, errs[0], "空です")
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		errs := ContactPersonNameValidate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Len(t, errs, 1)
		assert.EqualError(t, errs[0], "50文字より多いです")
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		errs := ContactPersonNameValidate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.Len(t, errs, 0)
	})
}
