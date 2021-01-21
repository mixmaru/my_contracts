package user

import (
	"github.com/mixmaru/my_contracts/core/domain/validators"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContactPersonNameValue_NewPresidentNameValue(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		name, err := NewPresidentNameValue("社長名")
		assert.NoError(t, err)
		assert.Equal(t, PresidentNameValue{"社長名"}, name)
	})

	t.Run("社長名が空文字だった時", func(t *testing.T) {
		name, err := NewPresidentNameValue("")
		assert.Error(t, err)
		assert.Equal(t, PresidentNameValue{}, name)
	})

	t.Run("社長名が50文字を超えていた時", func(t *testing.T) {
		name, err := NewPresidentNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Error(t, err)
		assert.Equal(t, PresidentNameValue{}, name)
	})

	t.Run("社長名が50文字だった時", func(t *testing.T) {
		name, err := NewPresidentNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Equal(t, PresidentNameValue{"0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789"}, name)
	})
}

func TestNameValue_PresidentNameValidate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		validErrs, err := PresidentNameValue{}.Validate("社長名")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		validErrs, err := PresidentNameValue{}.Validate("")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.EmptyStringValidError, validErrs[0])
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		validErrs, err := PresidentNameValue{}.Validate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.OverLengthStringValidError, validErrs[0])
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		validErrs, err := PresidentNameValue{}.Validate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})
}
