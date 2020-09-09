package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
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

func TestBillingCalculatorDomainService_ExecuteBilling(t *testing.T) {
	t.Run("渡した日時を実行日として_請求を実行する（billsとbill_detailsデータを作成する）", func(t *testing.T) {
		t.Run("2020/7/1を渡すと_7/1時点で使用権が開始されていて克つ_契約の課金開始日以降である使用権の使用量が請求される", func(t *testing.T) {
			db, err := db_connection.GetConnection()
			assert.NoError(t, err)

			productRep := repositories.NewProductRepository()
			contractRep := repositories.NewContractRepository()
			rightToUseRep := repositories.NewRightToUseRepository()
			userRep := repositories.NewUserRepository()

			////// 準備（2ユーザーに対して、6/1~6/30, 7/1~7/31, 8/1~8/31の未請求使用権データを作成する）
			product := createProduct("1000")
			user1, err := entities.NewUserIndividualEntity("ユーザー1")
			assert.NoError(t, err)
			user1Id, err := userRep.SaveUserIndividual(user1, db)

			productId, err := productRep.Save(product, db)
			assert.NoError(t, err)

			contract1 := entities.NewContractEntity(
				user1.Id(),
				productId,
				utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 6, 11, 0, 0, 0, 0),
			)
			contract1Id, err := contractRep.Create(contract1, db)
			assert.NoError(t, err)

			rightToUse1A := entities.NewRightToUseEntity(
				contract1Id,
				utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
			)
			rightToUse1AId, err := rightToUseRep.Create(rightToUse1A, db)
			assert.NoError(t, err)

			rightToUse1B := entities.NewRightToUseEntity(
				contract1Id,
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
			)
			rightToUse1BId, err := rightToUseRep.Create(rightToUse1B, db)
			assert.NoError(t, err)

			rightToUse1C := entities.NewRightToUseEntity(
				contract1Id,
				utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 9, 1, 0, 0, 0, 0),
			)
			rightToUse1CId, err := rightToUseRep.Create(rightToUse1C, db)
			assert.NoError(t, err)

			////// 実行
			ds := NewBillingCalculatorDomainService(
				productRep,
				contractRep,
				rightToUseRep,
			)
			err = ds.ExecuteBilling(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0))
			assert.NoError(t, err)

			////// 検証（billingデータを取得して検証する。2ユーザーの6/1~6/30, 7/1~7/31の請求分がbillsに作成される）
			billRep := repositories.NewBillRepository()
			actual1, err := billRep.GetById(user1Id, db)
			assert.NoError(t, err)

			assert.NotZero(t, actual1.Id())
			assert.Equal(t, utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0).String(), actual1.BillingDate().String())
			assert.Equal(t, user1Id, actual1.UserId())
			_, isNil, err := actual1.PaymentConfirmedAt()
			assert.Nil(t, isNil)
			total := actual1.TotalAmountExcludingTax()
			assert.Equal(t, decimal.NewFromInt(2000).String(), total.String())

			billDetails1 := actual1.BillDetails()
			assert.Len(t, billDetails1, 2)
			assert.NotZero(t, billDetails1[0].Id())
			assert.Equal(t, 1, billDetails1[0].OrderNum())
			assert.Equal(t, rightToUse1AId, billDetails1[0].RightToUseId())
			billingAmount := billDetails1[0].BillingAmount()
			assert.Equal(t, decimal.NewFromInt(1000).String(), billingAmount.String())

			assert.NotZero(t, billDetails1[1].Id())
			assert.Equal(t, 2, billDetails1[1].OrderNum())
			assert.Equal(t, rightToUse1BId, billDetails1[1].RightToUseId())
			billingAmount = billDetails1[1].BillingAmount()
			assert.Equal(t, decimal.NewFromInt(1000).String(), billingAmount.String())
			assert.Equal(t, 1, billDetails1[1].CreatedAt())
			assert.Equal(t, 1, billDetails1[1].UpdatedAt())

			//actual2, err := billRep.GetById(1, db)
			//assert.NoError(t, err)
			//
			//assert.NotZero(t, actual2.Id())
			//assert.Equal(t, utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0).String(), actual2.BillingDate().String())
			//assert.Equal(t, 1, actual2.UserId())
			//_, isNil, err = actual2.PaymentConfirmedAt()
			//assert.Nil(t, isNil)
			//
			//billDetails2 := actual2.BillDetails()
			//assert.Len(t, billDetails2, 2)
			//assert.NotZero(t, billDetails2[0].Id())
			//assert.Equal(t, 1, billDetails2[0].OrderNum())
			//assert.Equal(t, 1, billDetails2[0].RightToUseId())
			//assert.Equal(t, 1, billDetails2[0].BillingAmount())
			//assert.Equal(t, 1, billDetails2[0].CreatedAt())
			//assert.Equal(t, 1, billDetails2[0].UpdatedAt())
			//
			//assert.NotZero(t, billDetails2[1].Id())
			//assert.Equal(t, 1, billDetails2[1].OrderNum())
			//assert.Equal(t, 1, billDetails2[1].RightToUseId())
			//assert.Equal(t, 1, billDetails2[1].BillingAmount())
			//assert.Equal(t, 1, billDetails2[1].CreatedAt())
			//assert.Equal(t, 1, billDetails2[1].UpdatedAt())
		})
	})
}
