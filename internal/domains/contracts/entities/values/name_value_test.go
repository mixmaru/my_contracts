package values

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
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
		validErrs := NameValidate("個人顧客名")
		assert.Len(t, validErrs, 0)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		validErrs := NameValidate("")
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.EmptyStringValidError, validErrs[0])
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		validErrs := NameValidate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.OverLengthStringValidError, validErrs[0])
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		validErrs := NameValidate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.Len(t, validErrs, 0)
	})

}
