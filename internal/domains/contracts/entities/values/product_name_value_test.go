package values

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities/values/validators"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductNameValue_NewProductNameValue(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		name, err := NewProductNameValue("商品名")
		assert.NoError(t, err)
		assert.Equal(t, ProductNameValue{"商品名"}, name)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		name, err := NewProductNameValue("")
		assert.Error(t, err)
		assert.Equal(t, ProductNameValue{}, name)
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		name, err := NewProductNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Error(t, err)
		assert.Equal(t, ProductNameValue{}, name)
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		name, err := NewProductNameValue("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Equal(t, ProductNameValue{"0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789"}, name)
	})
}

func TestProductNameValue_ProductNameValidate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		validErrs := ProductNameValidate("商品名")
		assert.Len(t, validErrs, 0)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		validErrs := ProductNameValidate("")
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.EmptyStringValidError, validErrs[0])
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		validErrs := ProductNameValidate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.OverLengthStringValidError, validErrs[0])
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		validErrs := ProductNameValidate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.Len(t, validErrs, 0)
	})
}
