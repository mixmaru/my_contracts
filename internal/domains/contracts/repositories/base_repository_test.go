package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseRepository_selectOne(t *testing.T) {
	// テーブル事前削除
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	_, err = db.Exec("truncate table products cascade")
	assert.NoError(t, err)

	r := ProductRepository{}

	// 既存データ登録
	savedProductEntity, err := entities.NewProductEntity("商品名", "1000")
	assert.NoError(t, err)
	_, err = r.Save(savedProductEntity, nil)
	assert.NoError(t, err)

	t.Run("データがある時", func(t *testing.T) {
		// データ取得
		productRecord := data_mappers.ProductMapper{}
		productEntity := entities.ProductEntity{}
		noRow, err := selectOne(db, &productRecord, &productEntity, "select * from products where id =$1", savedProductEntity.Id())
		assert.NoError(t, err)
		assert.False(t, noRow)

		assert.Equal(t, savedProductEntity.Id(), productEntity.Id())
		assert.Equal(t, savedProductEntity.Name(), productEntity.Name())
		assert.Equal(t, savedProductEntity.Price(), productEntity.Price())
		assert.Equal(t, savedProductEntity.CreatedAt(), productEntity.CreatedAt())
		assert.Equal(t, savedProductEntity.UpdatedAt(), productEntity.UpdatedAt())
	})

	t.Run("データがない時", func(t *testing.T) {
		productRecord := data_mappers.ProductMapper{}
		productEntity := entities.ProductEntity{}
		noRow, err := selectOne(db, &productRecord, &productEntity, "select * from products where id =$1", -1000)
		assert.NoError(t, err)
		assert.True(t, noRow)
	})

	t.Run("渡すrecordとentityがアベコベだったとき", func(t *testing.T) {
		productRecord := data_mappers.ProductMapper{}
		userCorporationEntity := entities.UserCorporationEntity{}
		noRow, err := selectOne(db, &productRecord, &userCorporationEntity, "select * from products where id =$1", savedProductEntity.Id())
		assert.Error(t, err, "aaa")
		assert.True(t, noRow)
	})
}
