package product

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/tables"
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

func TestProductRepository_GetById(t *testing.T) {
	// テーブル事前削除
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	_, err = db.Exec("truncate table products cascade")
	assert.NoError(t, err)

	r := ProductRepository{}

	t.Run("データがある時", func(t *testing.T) {
		// データ登録
		productEntity := entities.NewProductEntity("商品名", decimal.NewFromFloat(1000))
		_, err = r.Save(productEntity, nil)
		assert.NoError(t, err)

		// データ取得
		loadedEntity, err := r.GetById(productEntity.Id(), nil)
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
		loadedEntity, err := r.GetById(-100, nil)
		assert.NoError(t, err)
		assert.Nil(t, loadedEntity)
	})
}

func TestProductRepository_selectOne(t *testing.T) {
	// テーブル事前削除
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	_, err = db.Exec("truncate table products cascade")
	assert.NoError(t, err)

	r := ProductRepository{}

	// 既存データ登録
	savedProductEntity := entities.NewProductEntity("商品名", decimal.NewFromFloat(1000))
	_, err = r.Save(savedProductEntity, nil)
	assert.NoError(t, err)

	t.Run("データがある時", func(t *testing.T) {
		// データ取得
		productRecord := tables.ProductRecord{}
		productEntity := entities.ProductEntity{}
		noRow, err := r.selectOne(db, &productRecord, &productEntity, "select * from products where id =$1", savedProductEntity.Id())
		assert.NoError(t, err)
		assert.False(t, noRow)

		assert.Equal(t, savedProductEntity.Id(), productEntity.Id())
		assert.Equal(t, savedProductEntity.Name(), productEntity.Name())
		assert.Equal(t, savedProductEntity.Price(), productEntity.Price())
		assert.Equal(t, savedProductEntity.CreatedAt(), productEntity.CreatedAt())
		assert.Equal(t, savedProductEntity.UpdatedAt(), productEntity.UpdatedAt())
	})

	t.Run("データがない時", func(t *testing.T) {
		productRecord := tables.ProductRecord{}
		productEntity := entities.ProductEntity{}
		noRow, err := r.selectOne(db, &productRecord, &productEntity, "select * from products where id =$1", -1000)
		assert.NoError(t, err)
		assert.True(t, noRow)
	})

	t.Run("渡すrecordとentityがアベコベだったとき", func(t *testing.T) {
		productRecord := tables.ProductRecord{}
		userCorporationEntity := entities.UserCorporationEntity{}
		noRow, err := r.selectOne(db, &productRecord, &userCorporationEntity, "select * from products where id =$1", savedProductEntity.Id())
		assert.Error(t, err, "aaa")
		assert.True(t, noRow)
	})
}
