package values

import (
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
		t.Skip()
		product, err := NewProductPriceValue("")
		assert.Error(t, err)
		assert.Equal(t, ProductPriceValue{}, product)
	})

	t.Run("priceがマイナスだったとき", func(t *testing.T) {
		t.Skip()
		product, err := NewProductPriceValue("-1000.01")
		assert.Error(t, err)
		assert.Equal(t, ProductPriceValue{}, product)
	})
}

func TestProductPriceValue_ProductPriceValidate(t *testing.T) {
	t.Skip()
	t.Run("正常系", func(t *testing.T) {
		errs := ProductPriceValidate("100.01")
		assert.Len(t, errs, 0)
	})

	t.Run("空文字だった時", func(t *testing.T) {
		errs := ProductPriceValidate("")
		assert.Len(t, errs, 1)
		assert.EqualError(t, errs[0], "空です")
	})

	t.Run("マイナスだったとき", func(t *testing.T) {
		errs := ProductPriceValidate("-1")
		assert.Len(t, errs, 1)
		assert.EqualError(t, errs[0], "マイナスです")
	})
}
