package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

func TestRightToUseRepository_Create(t *testing.T) {
	r := NewRightToUseRepository()
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	////// テスト用契約を作成する
	// 契約の作成
	savedContractId := createPreparedContractData(
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
		db,
	)

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

func TestRightToUseRepository_GetById(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// 事前に使用権を登録する
	r := NewRightToUseRepository()
	savedContractId := createPreparedContractData(
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
		db,
	)
	rightToUseEntity := entities.NewRightToUseEntity(
		savedContractId,
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
	)

	// 実行
	savedId, err := r.Create(rightToUseEntity, db)
	assert.NoError(t, err)

	t.Run("データがあればIdを渡すとデータが取得できる", func(t *testing.T) {
		// 実行
		actual, err := r.GetById(savedId, db)
		assert.NoError(t, err)

		// 検証
		assert.Equal(t, savedId, actual.Id())
		assert.Equal(t, savedContractId, actual.ContractId())
		assert.True(t, actual.ValidFrom().Equal(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)))
		assert.True(t, actual.ValidTo().Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
	})

	t.Run("データがなければidを渡すとnilが返る", func(t *testing.T) {
		// 実行
		actual, err := r.GetById(-100, db)
		assert.NoError(t, err)

		// 検証
		assert.Nil(t, actual)
	})
}

func TestRightToUseRepository_GetBillingTargetByBillingDate(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	r := NewRightToUseRepository()
	t.Run("請求実行日を渡すとその日以前で請求実行をしていない使用権データがuserId、rightToUserId順で全て返る", func(t *testing.T) {
		// 準備
		// 商品登録
		productId := createProduct(db)

		// user1
		// 6/1 ~ 6/30 の使用権（未請求）=> 取得される
		// 6/1 ~ 6/30 の使用権（請求済）
		// 7/1 ~ 7/31 の使用権（未請求）=> 取得される
		// 8/1 ~ 8/31 の使用権（未請求）
		userId1 := createUser(db)
		contractId1 := createContract(
			userId1,
			productId,
			utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
			utils.CreateJstTime(2020, 6, 2, 0, 0, 0, 0),
			db,
		)

		// 使用権
		rightToUseEntity1a := entities.NewRightToUseEntity(
			contractId1,
			utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
			utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
		)
		rightToUseId1a, err := r.Create(rightToUseEntity1a, db)
		assert.NoError(t, err)

		rightToUseEntity1b := entities.NewRightToUseEntity(
			contractId1,
			utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
			utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
		)
		rightToUseId1b, err := r.Create(rightToUseEntity1b, db)
		assert.NoError(t, err)

		rightToUseEntity1c := entities.NewRightToUseEntity(
			contractId1,
			utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
		)
		rightToUseId1c, err := r.Create(rightToUseEntity1c, db)
		assert.NoError(t, err)

		rightToUseEntity1d := entities.NewRightToUseEntity(
			contractId1,
			utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 9, 1, 0, 0, 0, 0),
		)
		rightToUseId1d, err := r.Create(rightToUseEntity1d, db)
		assert.NoError(t, err)

		// rightToUse1bに対して請求情報を登録しておく（請求済にしておく）
		billDetailEntity := entities.NewBillingDetailEntity(1, rightToUseId1b, decimal.NewFromInt(1000))
		billAgg := entities.NewBillingAggregation(utils.CreateJstTime(2020, 7, 1, 12, 0, 0, 0))
		err = billAgg.AddBillDetail(billDetailEntity)
		assert.NoError(t, err)

		billRep := NewBillRepository()
		billId, err := billRep.Create(billAgg, db)
		assert.NoError(t, err)

		assert.NotZero(t, rightToUseId1a)
		assert.NotZero(t, rightToUseId1b)
		assert.NotZero(t, rightToUseId1c)
		assert.NotZero(t, rightToUseId1d)
		assert.NotZero(t, billId)

		// user2
		// 6/1 ~ 6/30 の使用権（未請求）=> 取得される
		// 6/1 ~ 6/30 の使用権（請求済）
		// 7/1 ~ 7/31 の使用権（未請求）=> 取得される
		// 8/1 ~ 8/31 の使用権（未請求）

		// 実行
		billingDate := utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0)
		r.GetBillingTargetByBillingDate(billingDate)

		// 検証
	})
}

// 使用権データを作成するのに事前に必要なデータを準備する
func createPreparedContractData(contractDate, billingStartDate time.Time, executor gorp.SqlExecutor) int {
	// userの作成
	savedUserId := createUser(executor)
	// 商品の作成
	savedProductId := createProduct(executor)
	// 契約の作成
	savedContractId := createContract(
		savedUserId,
		savedProductId,
		contractDate,
		billingStartDate,
		executor,
	)
	return savedContractId
}

func createUser(executor gorp.SqlExecutor) int {
	userEntity, err := entities.NewUserIndividualEntity("個人太郎")
	if err != nil {
		panic("userEntity作成失敗")
	}
	userRepository := NewUserRepository()
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, executor)
	if err != nil {
		panic("userEntity保存失敗")
	}
	return savedUserId
}

func createProduct(executor gorp.SqlExecutor) int {
	productEntity, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), "1000")
	if err != nil {
		panic("productEntity作成失敗")
	}
	productRepository := NewProductRepository()
	savedProductId, err := productRepository.Save(productEntity, executor)
	if err != nil {
		panic("productEntity保存失敗")
	}
	return savedProductId
}

func createContract(userId, productId int, contractDate, billingStartDate time.Time, executor gorp.SqlExecutor) int {
	// 契約の作成
	contractEntity := entities.NewContractEntity(
		userId,
		productId,
		contractDate,
		billingStartDate,
	)
	contractRepository := NewContractRepository()
	savedContractId, err := contractRepository.Create(contractEntity, executor)
	if err != nil {
		panic("contractEntity保存失敗")
	}
	return savedContractId
}
