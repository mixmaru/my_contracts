package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// UserIndividualのインスタンス化をテスト
func TestContractEntity_NewContractEntity(t *testing.T) {
	// インスタンス化
	entity, err := NewContractEntity(1, 2)
	assert.NoError(t, err)

	// テスト
	assert.Equal(t, 1, entity.UserId())
	assert.Equal(t, 2, entity.ProductId())
}

//func TestContractEntity_NewProductEntityWithData(t *testing.T) {
//	// インスタンス化
//	createdAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
//	updatedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
//	productEntity, err := NewProductEntityWithData(1, "name", "1000", createdAt, updatedAt)
//	assert.NoError(t, err)
//
//	assert.Equal(t, 1, productEntity.Id())
//	assert.Equal(t, "name", productEntity.Name())
//	price := productEntity.Price()
//	assert.Equal(t, "1000", price.String())
//	assert.True(t, createdAt.Equal(productEntity.CreatedAt()))
//	assert.True(t, updatedAt.Equal(productEntity.UpdatedAt()))
//}
//
//func TestContractEntity_LoadData(t *testing.T) {
//	productEntity, err := NewProductEntity("name", "1000")
//	assert.NoError(t, err)
//
//	createdAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
//	updateAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
//	err = productEntity.LoadData(
//		1,
//		"name2",
//		"2000",
//		createdAt,
//		updateAt,
//	)
//	assert.NoError(t, err)
//
//	assert.Equal(t, 1, productEntity.Id())
//	assert.Equal(t, "name2", productEntity.Name())
//	price := productEntity.Price()
//	assert.True(t, price.Equal(decimal.NewFromFloat(2000)))
//	assert.Equal(t, createdAt, productEntity.CreatedAt())
//	assert.Equal(t, updateAt, productEntity.UpdatedAt())
//}
