package values

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values/validators"
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
		errs, err := ContactPersonNameValue{}.Validate("担当者名")
		assert.NoError(t, err)
		assert.Len(t, errs, 0)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		validErrs, err := ContactPersonNameValue{}.Validate("")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.EmptyStringValidError, validErrs[0])
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		validErrs, err := ContactPersonNameValue{}.Validate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.OverLengthStringValidError, validErrs[0])
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		validErrs, err := ContactPersonNameValue{}.Validate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})
}
