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
		validErrs, err := ProductPriceValidate("100.01")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 0)
	})

	t.Run("空文字だった時", func(t *testing.T) {
		validErrs, err := ProductPriceValidate("")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.EqualError(t, validErrs[0], "空です")
	})

	t.Run("マイナスだったとき", func(t *testing.T) {
		validErrs, err := ProductPriceValidate("-1")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.EqualError(t, validErrs[0], "マイナス値です")
	})

	t.Run("数値じゃなかったとき", func(t *testing.T) {
		validErrs, err := ProductPriceValidate("aaa")
		assert.NoError(t, err)
		assert.Len(t, validErrs, 1)
		assert.EqualError(t, validErrs[0], "数値ではありません")
	})
}
