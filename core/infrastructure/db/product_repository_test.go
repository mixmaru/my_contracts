package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductRepository_Save(t *testing.T) {
	t.Run("重複しない商品名を渡すと商品データを新規登録できる", func(t *testing.T) {
		// 準備
		db, err := GetConnection()
		assert.NoError(t, err)
		defer db.Db.Close()
		tran, err := db.Begin()
		assert.NoError(t, err)

		// 実行
		r := NewProductRepository()
		productEntity, err := product.NewProductEntity("商品", "1000")
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
	db, err := GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	r := NewProductRepository()

	// 検証用データ登録
	productEntity, err := product.NewProductEntity("商品", "1000")
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
	db, err := GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()
	_, err = db.Exec(`DELETE FROM product_price_monthlies ppm
WHERE ppm.product_id IN (
      SELECT id FROM products WHERE name = 'GetByNameテスト用商品'
)`)
	assert.NoError(t, err)
	_, err = db.Exec("DELETE FROM products WHERE name = 'GetByNameテスト用商品'")
	assert.NoError(t, err)

	r := NewProductRepository()

	// 検証用データ登録
	productEntity1, err := product.NewProductEntity("GetByNameテスト用商品", "1000")
	assert.NoError(t, err)
	savedId1, err := r.Save(productEntity1, db)
	assert.NoError(t, err)
	// 検証用データ登録
	productEntity2, err := product.NewProductEntity("GetByNameテスト用商品", "2000")
	assert.NoError(t, err)
	savedId2, err := r.Save(productEntity2, db)
	assert.NoError(t, err)

	t.Run("データがある時は商品名でデータを取得できる", func(t *testing.T) {
		// データ取得
		loadedEntities, err := r.GetByName("GetByNameテスト用商品", db)
		assert.NoError(t, err)

		// 検証
		assert.Len(t, loadedEntities, 2)
		// 1つめ
		assert.Equal(t, savedId1, loadedEntities[0].Id())
		assert.Equal(t, productEntity1.Name(), loadedEntities[0].Name())
		price, exist := loadedEntities[0].MonthlyPrice()
		assert.True(t, exist)
		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
		assert.NotZero(t, loadedEntities[0].CreatedAt())
		assert.NotZero(t, loadedEntities[0].UpdatedAt())
		// 2つめ
		assert.Equal(t, savedId2, loadedEntities[1].Id())
		assert.Equal(t, productEntity2.Name(), loadedEntities[1].Name())
		price, exist = loadedEntities[1].MonthlyPrice()
		assert.True(t, exist)
		assert.True(t, price.Equal(decimal.NewFromFloat(2000)))
		assert.NotZero(t, loadedEntities[1].CreatedAt())
		assert.NotZero(t, loadedEntities[1].UpdatedAt())
	})

	t.Run("データがない時は空スライスが返る", func(t *testing.T) {
		// データ取得
		loadedEntity, err := r.GetByName("", db) //空文字商品は登録できないようになってるので、カラ文字で検証する
		assert.NoError(t, err)
		assert.Len(t, loadedEntity, 0)
	})
}
func TestProductRepository_GetByRightToUseId(t *testing.T) {
	db, err := GetConnection()
	assert.NoError(t, err)

	r := NewProductRepository()

	t.Run("RightToUseIdを渡すと関連づいている商品データを返す", func(t *testing.T) {
		////// 準備
		test_db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		// テスト用userの登録
		userEntity, err := entities.NewUserIndividualEntity("請求計算用顧客")
		assert.NoError(t, err)
		userRep := repositories.NewUserRepository()
		savedUserId, err := userRep.SaveUserIndividual(userEntity, test_db)

		// 事前に31000円の商品を登録
		productEntity, err := product.NewProductEntity("商品", "3100")
		assert.NoError(t, err)
		productRep := NewProductRepository()
		savedProductId, err := productRep.Save(productEntity, db)
		assert.NoError(t, err)

		// 使用権を作成
		rightToUse := entities.NewRightToUseEntity(
			utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
		)
		// 契約を作成
		contractEntity := entities.NewContractEntity(
			savedUserId,
			savedProductId,
			utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
			utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
			[]*entities.RightToUseEntity{
				rightToUse,
			},
		)
		contractRep := repositories.NewContractRepository()
		savedContractId, err := contractRep.Create(contractEntity, test_db)
		assert.NoError(t, err)
		// 登録した契約データを再読込
		loadedContract, err := contractRep.GetById(savedContractId, test_db)

		////// 実行
		actual, err := r.GetByRightToUseId(loadedContract.RightToUses()[0].Id(), db)
		assert.NoError(t, err)

		// 検証
		assert.Equal(t, savedProductId, actual.Id())
		assert.Equal(t, productEntity.Name(), actual.Name())
		expectPrice, ok := productEntity.MonthlyPrice()
		assert.True(t, ok)
		actualPrice, ok := actual.MonthlyPrice()
		assert.True(t, ok)
		assert.Equal(t, expectPrice.String(), actualPrice.String())
		assert.NotZero(t, actual.CreatedAt())
		assert.NotZero(t, actual.UpdatedAt())
	})
}
