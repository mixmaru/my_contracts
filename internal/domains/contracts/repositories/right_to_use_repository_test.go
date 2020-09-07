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

func createBillTestData(db gorp.SqlExecutor) (rightToUseIds []int, billId int) {
	rightToUseRep := NewRightToUseRepository()

	// 商品登録
	productId := createProduct(db)

	// user作成
	// 6/1 ~ 6/30 の使用権（未請求）=> 取得される
	// 6/1 ~ 6/30 の使用権（請求済）
	// 7/1 ~ 7/31 の使用権（未請求）=> 取得される
	// 8/1 ~ 8/31 の使用権（未請求）
	userId := createUser(db)
	contractId1 := createContract(
		userId,
		productId,
		utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
		utils.CreateJstTime(2020, 6, 11, 0, 0, 0, 0),
		db,
	)

	// 使用権
	rightToUseEntity1 := entities.NewRightToUseEntity(
		contractId1,
		utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
		utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
	)
	rightToUseId1, err := rightToUseRep.Create(rightToUseEntity1, db)
	if err != nil {
		panic("rightToUseデータ保存失敗")
	}

	rightToUseEntity2 := entities.NewRightToUseEntity(
		contractId1,
		utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
		utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
	)
	rightToUseId2, err := rightToUseRep.Create(rightToUseEntity2, db)
	if err != nil {
		panic("rightToUseデータ保存失敗")
	}

	rightToUseEntity3 := entities.NewRightToUseEntity(
		contractId1,
		utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
	)
	rightToUseId3, err := rightToUseRep.Create(rightToUseEntity3, db)
	if err != nil {
		panic("rightToUseデータ保存失敗")
	}

	rightToUseEntity4 := entities.NewRightToUseEntity(
		contractId1,
		utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 9, 1, 0, 0, 0, 0),
	)
	rightToUseId4, err := rightToUseRep.Create(rightToUseEntity4, db)
	if err != nil {
		panic("rightToUseデータ保存失敗")
	}

	// rightToUse1bに対して請求情報を登録しておく（請求済にしておく）
	billDetailEntity := entities.NewBillingDetailEntity(1, rightToUseId2, decimal.NewFromInt(1000))
	billAgg := entities.NewBillingAggregation(utils.CreateJstTime(2020, 7, 1, 12, 0, 0, 0))
	err = billAgg.AddBillDetail(billDetailEntity)

	billRep := NewBillRepository()
	billId, err = billRep.Create(billAgg, db)
	if err != nil {
		panic("billデータ保存失敗")
	}

	return []int{rightToUseId1, rightToUseId2, rightToUseId3, rightToUseId4}, billId
}

func deleteTestRightToUseData(executor gorp.SqlExecutor) {
	query := `
DELETE FROM right_to_use WHERE id IN (
    SELECT rtu.id FROM right_to_use rtu
    LEFT OUTER JOIN bill_details bd on rtu.id = bd.right_to_use_id
    WHERE bd.id IS NULL
);
`
	_, err := executor.Exec(query)
	if err != nil {
		panic("事前データ削除失敗")
	}
}

func TestRightToUseRepository_GetBillingTargetByBillingDate(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)

	r := NewRightToUseRepository()
	t.Run("請求実行日を渡すとその日以前で請求実行をしていない使用権データがuserId、rightToUserId順で全て返る", func(t *testing.T) {
		// 準備
		// 事前に対象になる使用権を削除しておく
		tran, err := db.Begin()
		assert.NoError(t, err)

		deleteTestRightToUseData(tran)

		// 2userに対して、契約の課金開始日が6/11の以下の使用権を作成する。
		// 6/1 ~ 6/30 の使用権（未請求)
		// 6/1 ~ 6/30 の使用権（請求済）
		// 7/1 ~ 7/31 の使用権（未請求)
		// 8/1 ~ 8/31 の使用権（未請求）
		rightToUseIds1, _ := createBillTestData(tran)
		rightToUseIds2, _ := createBillTestData(tran)

		// 実行
		billingDate := utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)
		actual, err := r.GetBillingTargetByBillingDate(billingDate, tran)
		assert.NoError(t, err)

		err = tran.Commit()
		assert.NoError(t, err)

		// 検証
		// 6/1 ~ 6/30 の使用権（未請求) => 取得される（未請求だから）
		// 6/1 ~ 6/30 の使用権（請求済）=> 取得されない（請求実行済だから）
		// 7/1 ~ 7/31 の使用権（未請求) => 取得される（未請求だから）
		// 8/1 ~ 8/31 の使用権（未請求） => 取得されない（まだ請求しない使用権だから）
		assert.Len(t, actual, 4)
		assert.Equal(t, rightToUseIds1[0], actual[0].Id())
		assert.Equal(t, rightToUseIds1[2], actual[1].Id())
		assert.Equal(t, rightToUseIds2[0], actual[2].Id())
		assert.Equal(t, rightToUseIds2[2], actual[3].Id())
	})

	t.Run("渡した請求実行日が契約の課金開始日以前である使用権は返却データに含まれない_課金開始日以前の使用権には請求が発生しない", func(t *testing.T) {
		tran, err := db.Begin()
		assert.NoError(t, err)

		// 事前に影響のあるデータを削除
		deleteTestRightToUseData(tran)

		// 準備（以下の使用権データを作成する）
		// 6/1 ~ 6/30 の使用権（未請求）=> 取得されない（契約の課金開始日が6/11だから）
		// 6/1 ~ 6/30 の使用権（請求済）=> 取得されない
		// 7/1 ~ 7/31 の使用権（未請求）=> 取得されない
		// 8/1 ~ 8/31 の使用権（未請求）=> 取得されない
		_, _ = createBillTestData(tran)

		// 実行
		billingDate := utils.CreateJstTime(2020, 6, 10, 0, 0, 0, 0)
		actual, err := r.GetBillingTargetByBillingDate(billingDate, tran)
		assert.NoError(t, err)

		err = tran.Commit()
		assert.NoError(t, err)

		// 検証
		assert.Len(t, actual, 0)
	})

	t.Run("渡した請求実行日が契約の課金開始日ちょうどである使用権は返却データに含まれる", func(t *testing.T) {
		tran, err := db.Begin()
		assert.NoError(t, err)

		// 事前に影響のあるデータを削除
		deleteTestRightToUseData(tran)

		// 準備（以下の使用権データを作成する）
		// 6/1 ~ 6/30 の使用権（未請求）=> 取得される（契約の課金開始日が6/11だから）
		// 6/1 ~ 6/30 の使用権（請求済）=> 取得されない
		// 7/1 ~ 7/31 の使用権（未請求）=> 取得されない
		// 8/1 ~ 8/31 の使用権（未請求）=> 取得されない
		rightToUseIds, _ := createBillTestData(tran)

		// 実行
		billingDate := utils.CreateJstTime(2020, 6, 11, 0, 0, 0, 0)
		actual, err := r.GetBillingTargetByBillingDate(billingDate, tran)
		assert.NoError(t, err)

		err = tran.Commit()
		assert.NoError(t, err)

		// 検証
		assert.Len(t, actual, 1)
		assert.Equal(t, rightToUseIds[0], actual[0].Id())
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
