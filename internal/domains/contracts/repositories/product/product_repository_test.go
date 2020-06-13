package product

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductRepository_Save(t *testing.T) {
	// テーブル事前削除
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	_, err = db.Exec("truncate table products cascade")
	assert.NoError(t, err)

	r := ProductRepository{}
	productEntity := entities.NewProductEntity("商品名", decimal.NewFromFloat(1000))
	_, err = r.Save(productEntity, nil)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, productEntity.Id())
	assert.Equal(t, "商品名", productEntity.Name())
	price := productEntity.Price()
	assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
	assert.NotZero(t, productEntity.CreatedAt())
	assert.NotZero(t, productEntity.UpdatedAt())
}
