package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestRightToUseRepository_Create(t *testing.T) {
	r := NewRightToUseRepository()
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)

	////// テスト用契約を作成する
	// userの作成
	userEntity, err := entities.NewUserIndividualEntity("個人太郎")
	assert.NoError(t, err)
	userRepository := NewUserRepository()
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, db)
	assert.NoError(t, err)

	// 商品の作成
	// 重複しない商品名でテストを行う
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 10)
	name := "商品" + suffix
	productEntity, err := entities.NewProductEntity(name, "1000")
	assert.NoError(t, err)
	productRepository := NewProductRepository()
	savedProductId, err := productRepository.Save(productEntity, db)

	// 契約の作成
	contractEntity := entities.NewContractEntity(
		savedUserId,
		savedProductId,
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
	)
	contractRepository := NewContractRepository()
	savedContractId, err := contractRepository.Create(contractEntity, db)
	assert.NoError(t, err)

	t.Run("権利エンティティとdbコネクションを渡すとDBへ新規保存されて_保存Idを返す", func(t *testing.T) {
		// 準備
		entity := entities.NewRightToUseEntity(
			savedContractId,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
		)

		// 実行
		savedId, err := r.Create(entity, db)
		assert.NoError(t, err)

		//検証
		assert.NotZero(t, savedId)
	})
}

//
//import (
//	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
//	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
//	"github.com/mixmaru/my_contracts/internal/lib/decimal"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestProductRepository_Save(t *testing.T) {
//	// テーブル事前削除
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//	_, err = db.Exec("truncate table products cascade")
//	assert.NoError(t, err)
//
//	tran, err := db.Begin()
//	assert.NoError(t, err)
//
//	r := NewProductRepository()
//	productEntity, err := entities.NewProductEntity("商品名", "1000")
//	assert.NoError(t, err)
//	savedId, err := r.Save(productEntity, tran)
//	assert.NoError(t, err)
//	err = tran.Commit()
//	assert.NoError(t, err)
//	assert.NotEqual(t, 0, savedId)
//}
//
//func TestProductRepository_GetById(t *testing.T) {
//	// テーブル事前削除
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//	_, err = db.Exec("truncate table products cascade")
//	assert.NoError(t, err)
//
//	r := NewProductRepository()
//
//	// 検証用データ登録
//	productEntity, err := entities.NewProductEntity("商品名", "1000")
//	assert.NoError(t, err)
//	savedId, err := r.Save(productEntity, db)
//	assert.NoError(t, err)
//
//	t.Run("データがある時", func(t *testing.T) {
//		// データ取得
//		loadedEntity, err := r.GetById(savedId, db)
//		assert.NoError(t, err)
//
//		assert.Equal(t, savedId, loadedEntity.Id())
//		assert.Equal(t, "商品名", loadedEntity.Name())
//		price, exist := loadedEntity.MonthlyPrice()
//		assert.True(t, exist)
//		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
//		assert.NotZero(t, loadedEntity.CreatedAt())
//		assert.NotZero(t, loadedEntity.UpdatedAt())
//	})
//
//	t.Run("データがない時", func(t *testing.T) {
//		// データ取得
//		loadedEntity, err := r.GetById(-100, db)
//		assert.NoError(t, err)
//		assert.Nil(t, loadedEntity)
//	})
//}
//
//func TestProductRepository_GetByName(t *testing.T) {
//	// テーブル事前削除
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//	_, err = db.Exec("truncate table products cascade")
//	assert.NoError(t, err)
//
//	r := NewProductRepository()
//
//	// 検証用データ登録
//	productEntity, err := entities.NewProductEntity("商品名", "1000")
//	assert.NoError(t, err)
//	savedId, err := r.Save(productEntity, db)
//	assert.NoError(t, err)
//
//	t.Run("データがある時", func(t *testing.T) {
//
//		// データ取得
//		loadedEntity, err := r.GetByName("商品名", db)
//		assert.NoError(t, err)
//
//		assert.Equal(t, savedId, loadedEntity.Id())
//		assert.Equal(t, "商品名", loadedEntity.Name())
//		price, exist := loadedEntity.MonthlyPrice()
//		assert.True(t, exist)
//		assert.True(t, price.Equal(decimal.NewFromFloat(1000)))
//		assert.NotZero(t, loadedEntity.CreatedAt())
//		assert.NotZero(t, loadedEntity.UpdatedAt())
//	})
//
//	t.Run("データがない時", func(t *testing.T) {
//		// データ取得
//		loadedEntity, err := r.GetByName("存在しない商品", db)
//		assert.NoError(t, err)
//		assert.Nil(t, loadedEntity)
//	})
//}
