package product

import (
	"github.com/mixmaru/my_contracts/core/domain/validators"
	"testing"

	"github.com/stretchr/testify/assert"
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
		validErrs, err := ProductNameValue{}.Validate("商品名")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})

	t.Run("名前が空文字だった時", func(t *testing.T) {
		validErrs, err := ProductNameValue{}.Validate("")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.EmptyStringValidError, validErrs[0])
	})

	t.Run("名前が50文字を超えていた時", func(t *testing.T) {
		validErrs, err := ProductNameValue{}.Validate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789a")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.OverLengthStringValidError, validErrs[0])
	})

	t.Run("名前が50文字だった時", func(t *testing.T) {
		validErrs, err := ProductNameValue{}.Validate("0123456789０１２３４５６７８９0123456789０１２３４５６７８９0123456789")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})
}
