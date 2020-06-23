package entities

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserIndividualのインスタンス化をテスト
func TestProductEntity_NewProductEntity(t *testing.T) {
	// インスタンス化
	productEntity := NewProductEntity("name", decimal.NewFromFloat(1000))

	// テスト
	assert.Equal(t, "name", productEntity.Name())
	price := productEntity.Price()
	assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
}
func TestProductEntity_NewProductEntityWithData(t *testing.T) {
	// インスタンス化
	price := decimal.NewFromFloat(1000)
	createdAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	productEntity := NewProductEntityWithData(1, "name", price, createdAt, updatedAt)

	assert.Equal(t, 1, productEntity.Id())
	assert.Equal(t, "name", productEntity.Name())
	assert.True(t, price.Equal(productEntity.Price()))
	assert.True(t, createdAt.Equal(productEntity.CreatedAt()))
	assert.True(t, updatedAt.Equal(productEntity.UpdatedAt()))
}

func TestProductEntity_LoadData(t *testing.T) {
	productEntity := NewProductEntity("name", decimal.NewFromFloat(1000))
	createdAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	updateAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	productEntity.LoadData(
		1,
		"name2",
		decimal.NewFromFloat(2000),
		createdAt,
		updateAt,
	)

	assert.Equal(t, 1, productEntity.Id())
	assert.Equal(t, "name2", productEntity.Name())
	price := productEntity.Price()
	assert.True(t, price.Equal(decimal.NewFromFloat(2000)))
	assert.Equal(t, createdAt, productEntity.CreatedAt())
	assert.Equal(t, updateAt, productEntity.UpdatedAt())
}
