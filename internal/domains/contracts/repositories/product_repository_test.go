package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductRepository_Save(t *testing.T) {
	t.Run("重複しない商品名を渡すと商品データを新規登録できる", func(t *testing.T) {
		// 準備
		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		defer db.Db.Close()
		tran, err := db.Begin()
		assert.NoError(t, err)

		// 実行
		r := NewProductRepository()
		productEntity, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), "1000")
		assert.NoError(t, err)
		savedId, err := r.Save(productEntity, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		// 検証
		assert.NotEqual(t, 0, savedId)
	})
}

func TestProductRepository_GetById(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	r := NewProductRepository()

	// 検証用データ登録
	productEntity, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), "1000")
	assert.NoError(t, err)
	savedId, err := r.Save(productEntity, db)
	assert.NoError(t, err)

	t.Run("データがある時はIdで取得できる", func(t *testing.T) {
		// データ取得
		loadedEntity, err := r.GetById(savedId, db)
		assert.NoError(t, err)

		// 検証
		assert.Equal(t, savedId, loadedEntity.Id())
		assert.Equal(t, productEntity.Name(), loadedEntity.Name())
		price, exist := loadedEntity.MonthlyPrice()
		assert.True(t, exist)
		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
		assert.NotZero(t, loadedEntity.CreatedAt())
		assert.NotZero(t, loadedEntity.UpdatedAt())
	})

	t.Run("データがない時はnilが返ってくる", func(t *testing.T) {
		// データ取得
		loadedEntity, err := r.GetById(-100, db)
		assert.NoError(t, err)
		assert.Nil(t, loadedEntity)
	})
}

func TestProductRepository_GetByName(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	//defer db.Db.Close()
	//_, err = db.Exec("truncate table products cascade")
	//assert.NoError(t, err)

	r := NewProductRepository()

	// 検証用データ登録
	productEntity, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), "1000")
	assert.NoError(t, err)
	savedId, err := r.Save(productEntity, db)
	assert.NoError(t, err)

	t.Run("データがある時は商品名でデータを取得できる", func(t *testing.T) {
		// データ取得
		loadedEntity, err := r.GetByName(productEntity.Name(), db)
		assert.NoError(t, err)

		// 検証
		assert.Equal(t, savedId, loadedEntity.Id())
		assert.Equal(t, productEntity.Name(), loadedEntity.Name())
		price, exist := loadedEntity.MonthlyPrice()
		assert.True(t, exist)
		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
		assert.NotZero(t, loadedEntity.CreatedAt())
		assert.NotZero(t, loadedEntity.UpdatedAt())
	})

	t.Run("データがない時はnilが返る", func(t *testing.T) {
		// データ取得
		loadedEntity, err := r.GetByName("", db) //空文字商品は登録できないようになってるので、カラ文字で検証する
		assert.NoError(t, err)
		assert.Nil(t, loadedEntity)
	})
}
