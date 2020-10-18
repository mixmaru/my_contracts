package product

import (
	"testing"
	"time"

	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/stretchr/testify/assert"
)

// UserIndividualのインスタンス化をテスト
func TestProductEntity_NewProductEntity(t *testing.T) {
	t.Run("名前と価格を渡すと_今は暫定的に_月額1000円商品としてインスタンス化される", func(t *testing.T) {
		// インスタンス化
		productEntity, err := NewProductEntity("name", "1000")
		assert.NoError(t, err)

		// テスト
		assert.Equal(t, "name", productEntity.Name())
		price, exist := productEntity.MonthlyPrice()
		assert.True(t, exist)
		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
	})
}
func TestProductEntity_NewProductEntityWithData(t *testing.T) {
	// インスタンス化
	createdAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	productEntity, err := NewProductEntityWithData(1, "name", "1000", createdAt, updatedAt)
	assert.NoError(t, err)

	assert.Equal(t, 1, productEntity.Id())
	assert.Equal(t, "name", productEntity.Name())
	price, exist := productEntity.MonthlyPrice()
	assert.True(t, exist)
	assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
	assert.True(t, createdAt.Equal(productEntity.CreatedAt()))
	assert.True(t, updatedAt.Equal(productEntity.UpdatedAt()))
}

func TestProductEntity_LoadData(t *testing.T) {
	productEntity, err := NewProductEntity("name", "1000")
	assert.NoError(t, err)

	createdAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	updateAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	err = productEntity.LoadData(
		1,
		"name2",
		"2000",
		createdAt,
		updateAt,
	)
	assert.NoError(t, err)

	assert.Equal(t, 1, productEntity.Id())
	assert.Equal(t, "name2", productEntity.Name())
	price, exist := productEntity.MonthlyPrice()
	assert.True(t, exist)
	assert.True(t, price.Equal(decimal.NewFromFloat(2000)))
	assert.Equal(t, createdAt, productEntity.CreatedAt())
	assert.Equal(t, updateAt, productEntity.UpdatedAt())
}

func TestProductEntity_MonthlyPrice(t *testing.T) {
	t.Run("月契約が存在する商品なら_月額定価が返ってくる", func(t *testing.T) {
		productEntity, err := NewProductEntity("name", "1000")
		assert.NoError(t, err)

		price, exist := productEntity.MonthlyPrice()
		assert.True(t, exist)
		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
	})
	t.Run("月契約が存在しない商品なら_errorが返ってくる", func(t *testing.T) {
		t.Skip("まだ月契約しか設定できないのでスキップ")
	})
}

func TestProductEntity_GetTermType(t *testing.T) {
	t.Run("月契約だったらmonthlyが返る", func(t *testing.T) {
		productEntity, err := NewProductEntity("name", "1000")
		assert.NoError(t, err)

		expect, err := productEntity.GetTermType()
		assert.NoError(t, err)
		assert.Equal(t, TermMonthly, expect)
	})

}
