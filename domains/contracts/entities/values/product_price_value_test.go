package values

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities/values/validators"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductPriceValue_NewNameValue(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		product, err := NewProductPriceValue("1000.01")
		assert.NoError(t, err)

		value := product.Value()
		assert.Equal(t, "1000.01", value.String())
	})

	t.Run("priceが空文字だった時", func(t *testing.T) {
		product, err := NewProductPriceValue("")
		assert.Error(t, err)
		assert.Equal(t, ProductPriceValue{}, product)
	})

	t.Run("priceがマイナスだったとき", func(t *testing.T) {
		product, err := NewProductPriceValue("-1000.01")
		assert.Error(t, err)
		assert.Equal(t, ProductPriceValue{}, product)
	})
}

func TestProductPriceValue_ProductPriceValidate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		validErrs, err := ProductPriceValue{}.Validate("100.01")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})

	t.Run("空文字だった時", func(t *testing.T) {
		validErrs, err := ProductPriceValue{}.Validate("")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.EmptyStringValidError, validErrs[0])
	})

	t.Run("マイナスだったとき", func(t *testing.T) {
		validErrs, err := ProductPriceValue{}.Validate("-1")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.NegativeValidError, validErrs[0])
	})

	t.Run("数値じゃなかったとき", func(t *testing.T) {
		validErrs, err := ProductPriceValue{}.Validate("aaa")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.Equal(t, validators.NumericStringValidError, validErrs[0])
	})

	t.Run("文字列じゃなかったとき", func(t *testing.T) {
		validErrs, err := ProductPriceValue{}.Validate(1111)
		assert.Error(t, err)
		assert.Len(t, validErrs, 0)
	})
}
