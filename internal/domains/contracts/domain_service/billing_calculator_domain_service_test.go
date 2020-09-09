package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
)

func createProduct(price string) *entities.ProductEntity {
	db, err := db_connection.GetConnection()
	if err != nil {
		panic("db接続エラー")
	}
	productEntity, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), price)
	if err != nil {
		panic("データ作成エラー")
	}
	rep := repositories.NewProductRepository()
	id, err := rep.Save(productEntity, db)
	if err != nil {
		panic("データ保存エラー")
	}
	entity, err := rep.GetById(id, db)
	if err != nil {
		panic("データ取得エラー")
	}
	return entity
}

func TestBillingCalculatorDomainService_BillingAmount(t *testing.T) {
	// db接続作成
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)

	////// テスト用のデータの作成
	// テスト用userの登録
	userEntity, err := entities.NewUserIndividualEntity("請求計算用顧客")
	assert.NoError(t, err)
	userRep := repositories.NewUserRepository()
	savedUserId, err := userRep.SaveUserIndividual(userEntity, db)

	// 事前に31000円の商品を登録
	productEntity := createProduct("31000")
	savedProductId := productEntity.Id()
	assert.NoError(t, err)

	// ドメインサービスインスタンス化
	billingDS := NewBillingCalculatorDomainService(repositories.NewProductRepository(), repositories.NewContractRepository(), repositories.NewRightToUseRepository())

	t.Run("31日ある月", func(t *testing.T) {
		t.Run("使用権の有効期間が契約の課金単位期間_月払いとか年払いとか_の満了期間だと満額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成と課金開始日が同日）
			contractEntity := entities.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
			)
			contractRep := repositories.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, db)
			assert.NoError(t, err)

			// 使用権を作成
			rightToUse := entities.NewRightToUseEntity(
				savedContractId,
				utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
			)
			rightToUseRep := repositories.NewRightToUseRepository()
			savedRightToUseId, err := rightToUseRep.Create(rightToUse, db)
			assert.NoError(t, err)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedRightToUseId, db)
			assert.NoError(t, err)

			// 検証
			assert.Equal(t, "31000", billingAmount.String())
		})

		t.Run("使用権の有効期間に契約の課金開始日が食い込んでいると_課金開始日前は含まれず日割り金額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成から課金開始日まで5日間ある）
			contractEntity := entities.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 1, 6, 0, 0, 0, 0),
			)
			contractRep := repositories.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, db)
			assert.NoError(t, err)

			// 使用権を作成
			rightToUse := entities.NewRightToUseEntity(
				savedContractId,
				utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
			)
			rightToUseRep := repositories.NewRightToUseRepository()

			savedRightToUseId, err := rightToUseRep.Create(rightToUse, db)
			assert.NoError(t, err)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedRightToUseId, db)
			assert.NoError(t, err)

			// 検証
			assert.Equal(t, "26000", billingAmount.String())
		})
	})

	t.Run("30日ある月", func(t *testing.T) {
		t.Run("使用権の有効期間が契約の課金単位期間_月払いとか年払いとか_の満了期間だと満額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成と課金開始日が同日）
			contractEntity := entities.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
			)
			contractRep := repositories.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, db)
			assert.NoError(t, err)

			// 使用権を作成
			rightToUse := entities.NewRightToUseEntity(
				savedContractId,
				utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
			)
			rightToUseRep := repositories.NewRightToUseRepository()
			savedRightToUseId, err := rightToUseRep.Create(rightToUse, db)
			assert.NoError(t, err)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedRightToUseId, db)
			assert.NoError(t, err)

			// 検証
			assert.Equal(t, "31000", billingAmount.String())
		})

		t.Run("使用権の期間が1ヶ月に足りず満了期間ではない場合日割り金額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成から課金開始日まで5日間ある）
			contractEntity := entities.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 4, 6, 0, 0, 0, 0),
			)
			contractRep := repositories.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, db)
			assert.NoError(t, err)

			// 使用権を作成
			rightToUse := entities.NewRightToUseEntity(
				savedContractId,
				utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
			)
			rightToUseRep := repositories.NewRightToUseRepository()

			savedRightToUseId, err := rightToUseRep.Create(rightToUse, db)
			assert.NoError(t, err)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedRightToUseId, db)
			assert.NoError(t, err)

			// 検証（日割り金額になる 31000 / 30 * 25）で端数切り捨て
			assert.Equal(t, "25833", billingAmount.String())
		})
	})

	t.Run("29日ある月", func(t *testing.T) {
		t.Run("使用権の有効期間が契約の課金単位期間_月払いとか年払いとか_の満了期間だと満額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成と課金開始日が同日）
			contractEntity := entities.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
			)
			contractRep := repositories.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, db)
			assert.NoError(t, err)

			// 使用権を作成
			rightToUse := entities.NewRightToUseEntity(
				savedContractId,
				utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
			)
			rightToUseRep := repositories.NewRightToUseRepository()
			savedRightToUseId, err := rightToUseRep.Create(rightToUse, db)
			assert.NoError(t, err)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedRightToUseId, db)
			assert.NoError(t, err)

			// 検証
			assert.Equal(t, "31000", billingAmount.String())
		})

		t.Run("使用権の期間が1ヶ月に足りず満了期間ではない場合日割り金額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成から課金開始日まで5日間ある）
			contractEntity := entities.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 2, 6, 0, 0, 0, 0),
			)
			contractRep := repositories.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, db)
			assert.NoError(t, err)

			// 使用権を作成
			rightToUse := entities.NewRightToUseEntity(
				savedContractId,
				utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
			)
			rightToUseRep := repositories.NewRightToUseRepository()

			savedRightToUseId, err := rightToUseRep.Create(rightToUse, db)
			assert.NoError(t, err)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedRightToUseId, db)
			assert.NoError(t, err)

			// 検証（日割り金額になる 31000 / 29 * 24）で端数切り捨て
			assert.Equal(t, "25655", billingAmount.String())
		})
	})
}

func createTestData(executor gorp.SqlExecutor, t *testing.T) (userId, rightToUse1Id, rightToUse2Id, rightToUse3Id int) {
	contractRep := repositories.NewContractRepository()
	rightToUseRep := repositories.NewRightToUseRepository()
	userRep := repositories.NewUserRepository()

	////// 準備（2ユーザーに対して、6/1~6/30, 7/1~7/31, 8/1~8/31の未請求使用権データを作成する）
	product := createProduct("1000")
	user1, err := entities.NewUserIndividualEntity("ユーザー1")
	assert.NoError(t, err)
	user1Id, err := userRep.SaveUserIndividual(user1, executor)

	contract1 := entities.NewContractEntity(
		user1Id,
		product.Id(),
		utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 6, 11, 0, 0, 0, 0),
	)
	contract1Id, err := contractRep.Create(contract1, executor)
	assert.NoError(t, err)

	rightToUse1A := entities.NewRightToUseEntity(
		contract1Id,
		utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
	)
	rightToUse1AId, err := rightToUseRep.Create(rightToUse1A, executor)
	assert.NoError(t, err)

	rightToUse1B := entities.NewRightToUseEntity(
		contract1Id,
		utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
	)
	rightToUse1BId, err := rightToUseRep.Create(rightToUse1B, executor)
	assert.NoError(t, err)

	rightToUse1C := entities.NewRightToUseEntity(
		contract1Id,
		utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 9, 1, 0, 0, 0, 0),
	)
	rightToUse1CId, err := rightToUseRep.Create(rightToUse1C, executor)
	assert.NoError(t, err)

	return user1Id, rightToUse1AId, rightToUse1BId, rightToUse1CId
}

func assertBill(t *testing.T, actual, expect *entities.BillAggregation) {
	assert.NotZero(t, actual.Id())
	assert.Equal(t, utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0).String(), actual.BillingDate().String())
	assert.Equal(t, expect.UserId(), actual.UserId())
	_, isNil, err := actual.PaymentConfirmedAt()
	assert.NoError(t, err)
	assert.Nil(t, isNil)
	actualTotal := actual.TotalAmountExcludingTax()
	expectTotal := expect.TotalAmountExcludingTax()
	assert.Equal(t, expectTotal.String(), actualTotal.String())

	actualDetails := actual.BillDetails()
	expectDetails := expect.BillDetails()
	assert.Equal(t, len(expectDetails), len(actualDetails))
	assert.NotZero(t, actualDetails[0].Id())
	assert.Equal(t, expectDetails[0].OrderNum(), actualDetails[0].OrderNum())
	assert.Equal(t, expectDetails[0].RightToUseId(), actualDetails[0].RightToUseId())
	actualBillingAmount := actualDetails[0].BillingAmount()
	expectBillingAmount := expectDetails[0].BillingAmount()
	assert.Equal(t, expectBillingAmount.String(), actualBillingAmount.String())

	assert.NotZero(t, actualDetails[1].Id())
	assert.Equal(t, expectDetails[1].OrderNum(), actualDetails[1].OrderNum())
	assert.Equal(t, expectDetails[1].RightToUseId(), actualDetails[1].RightToUseId())
	actualBillingAmount = actualDetails[1].BillingAmount()
	expectBillingAmount = expectDetails[1].BillingAmount()
	assert.Equal(t, expectBillingAmount.String(), actualBillingAmount.String())
}

func TestBillingCalculatorDomainService_ExecuteBilling(t *testing.T) {
	t.Run("渡した日時を実行日として_請求を実行する（billsとbill_detailsデータを作成する）", func(t *testing.T) {
		t.Run("2020/7/1を渡すと_7/1時点で使用権が開始されていて克つ_契約の課金開始日以降である使用権の使用量が請求される", func(t *testing.T) {
			db, err := db_connection.GetConnection()
			assert.NoError(t, err)
			tran, err := db.Begin()

			// テストデータ作成
			user1Id, rightToUse1AId, rightToUse1BId, _ := createTestData(tran, t)

			////// 実行
			ds := NewBillingCalculatorDomainService(
				repositories.NewProductRepository(),
				repositories.NewContractRepository(),
				repositories.NewRightToUseRepository(),
			)
			err = ds.ExecuteBilling(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0))
			assert.NoError(t, err)

			////// 検証（billingデータを取得して検証する。2ユーザーの6/1~6/30, 7/1~7/31の請求分がbillsに作成される）

			billRep := repositories.NewBillRepository()
			actual1, err := billRep.GetById(user1Id, tran)
			assert.NoError(t, err)
			err = tran.Commit()
			assert.NoError(t, err)

			expect := entities.NewBillingAggregation(
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
				user1Id,
			)
			err = expect.AddBillDetail(entities.NewBillingDetailEntity(1, rightToUse1AId, decimal.NewFromInt(1000)))
			assert.NoError(t, err)
			err = expect.AddBillDetail(entities.NewBillingDetailEntity(2, rightToUse1BId, decimal.NewFromInt(1000)))
			assert.NoError(t, err)

			assertBill(t, actual1, expect)
		})
	})
}
