package repositories

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
	defer db.Db.Close()
	_, err = db.Exec("truncate table products cascade")
	assert.NoError(t, err)

	r := NewProductRepository()
	productEntity, err := entities.NewProductEntity("商品名", "1000")
	assert.NoError(t, err)
	savedId, err := r.Save(productEntity, db)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, savedId)
}

func TestProductRepository_GetById(t *testing.T) {
	// テーブル事前削除
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()
	_, err = db.Exec("truncate table products cascade")
	assert.NoError(t, err)

	r := NewProductRepository()

	t.Run("データがある時", func(t *testing.T) {
		// データ登録
		productEntity, err := entities.NewProductEntity("商品名", "1000")
		assert.NoError(t, err)
		_, err = r.Save(productEntity, db)
		assert.NoError(t, err)

		// データ取得
		loadedEntity, err := r.GetById(productEntity.Id(), db)
		assert.NoError(t, err)

		assert.Equal(t, productEntity.Id(), loadedEntity.Id())
		assert.Equal(t, "商品名", loadedEntity.Name())
		price := loadedEntity.Price()
		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
		assert.True(t, loadedEntity.CreatedAt().Equal(productEntity.CreatedAt()))
		assert.True(t, loadedEntity.UpdatedAt().Equal(productEntity.UpdatedAt()))
	})

	t.Run("データがない時", func(t *testing.T) {
		// データ取得
		loadedEntity, err := r.GetById(-100, db)
		assert.NoError(t, err)
		assert.Nil(t, loadedEntity)
	})
}

func TestProductRepository_GetByName(t *testing.T) {
	// テーブル事前削除
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()
	_, err = db.Exec("truncate table products cascade")
	assert.NoError(t, err)

	r := NewProductRepository()

	t.Run("データがある時", func(t *testing.T) {
		// データ登録
		productEntity, err := entities.NewProductEntity("商品名", "1000")
		assert.NoError(t, err)
		_, err = r.Save(productEntity, db)
		assert.NoError(t, err)

		// データ取得
		loadedEntity, err := r.GetByName("商品名", db)
		assert.NoError(t, err)

		assert.Equal(t, productEntity.Id(), loadedEntity.Id())
		assert.Equal(t, "商品名", loadedEntity.Name())
		price := loadedEntity.Price()
		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
		assert.True(t, loadedEntity.CreatedAt().Equal(productEntity.CreatedAt()))
		assert.True(t, loadedEntity.UpdatedAt().Equal(productEntity.UpdatedAt()))
	})

	t.Run("データがない時", func(t *testing.T) {
		// データ取得
		loadedEntity, err := r.GetByName("存在しない商品", db)
		assert.NoError(t, err)
		assert.Nil(t, loadedEntity)
	})
}
