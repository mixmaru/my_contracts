package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// UserIndividualのインスタンス化をテスト
func TestContractEntity_NewContractEntity(t *testing.T) {
	// インスタンス化
	entity := NewContractEntity(1, 2)

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

func TestContractEntity_LoadData(t *testing.T) {
	contractEntity, err := NewContractEntity(100, 200)
	assert.NoError(t, err)

	createdAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	updateAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	err = contractEntity.LoadData(
		1,
		2,
		3,
		createdAt,
		updateAt,
	)
	assert.NoError(t, err)

	assert.Equal(t, 1, contractEntity.Id())
	assert.Equal(t, 2, contractEntity.UserId())
	assert.Equal(t, 3, contractEntity.ProductId())
	assert.Equal(t, createdAt, contractEntity.CreatedAt())
	assert.Equal(t, updateAt, contractEntity.UpdatedAt())
}
