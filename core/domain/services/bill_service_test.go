package services

import (
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBillService_BillingAmount(t *testing.T) {
	// db接続作成
	conn, err := db.GetConnection()
	assert.NoError(t, err)

	////// テスト用のデータの作成
	// テスト用userの登録
	userEntity, err := user.NewUserIndividualEntity("請求計算用顧客")
	assert.NoError(t, err)
	userRep := db.NewUserRepository()
	savedUserId, err := userRep.SaveUserIndividual(userEntity, conn)

	// 事前に31000円の商品を登録
	productEntity := createProduct("31000")
	savedProductId := productEntity.Id()
	assert.NoError(t, err)

	// ドメインサービスインスタンス化
	billingDS := NewBillService(
		db.NewProductRepository(),
		db.NewContractRepository(),
		db.NewBillRepository(),
	)

	t.Run("31日ある月", func(t *testing.T) {
		t.Run("使用権の有効期間が契約の課金単位期間_月払いとか年払いとか_の満了期間だと満額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成と課金開始日が同日）
			contractEntity := contract.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
				[]*contract.RightToUseEntity{
					contract.NewRightToUseEntity(
						utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
						utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
					),
				},
			)
			contractRep := db.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, conn)
			assert.NoError(t, err)

			// リロード
			savedContract, err := contractRep.GetById(savedContractId, conn)
			assert.NoError(t, err)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedContract.RightToUses()[0], savedContract.BillingStartDate(), conn)
			assert.NoError(t, err)

			// 検証
			assert.Equal(t, "31000", billingAmount.String())
		})

		t.Run("使用権の有効期間に契約の課金開始日が食い込んでいると_課金開始日前は含まれず日割り金額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成から課金開始日まで5日間ある）
			contractEntity := contract.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 1, 6, 0, 0, 0, 0),
				[]*contract.RightToUseEntity{
					contract.NewRightToUseEntity(
						utils.CreateJstTime(2020, 1, 1, 15, 11, 36, 123456),
						utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
					),
				},
			)
			contractRep := db.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, conn)
			assert.NoError(t, err)
			// リロード
			savedContract, err := contractRep.GetById(savedContractId, conn)
			assert.NoError(t, err)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedContract.RightToUses()[0], savedContract.BillingStartDate(), conn)
			assert.NoError(t, err)

			// 検証
			assert.Equal(t, "26000", billingAmount.String())
		})
	})

	t.Run("30日ある月", func(t *testing.T) {
		t.Run("使用権の有効期間が契約の課金単位期間_月払いとか年払いとか_の満了期間だと満額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成と課金開始日が同日）
			contractEntity := contract.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
				[]*contract.RightToUseEntity{
					contract.NewRightToUseEntity(
						utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
						utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					),
				},
			)
			contractRep := db.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, conn)
			assert.NoError(t, err)
			// リロード
			savedContract, err := contractRep.GetById(savedContractId, conn)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedContract.RightToUses()[0], savedContract.BillingStartDate(), conn)
			assert.NoError(t, err)

			// 検証
			assert.Equal(t, "31000", billingAmount.String())
		})

		t.Run("使用権の期間が1ヶ月に足りず満了期間ではない場合日割り金額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成から課金開始日まで5日間ある）
			contractEntity := contract.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 4, 6, 0, 0, 0, 0),
				[]*contract.RightToUseEntity{
					contract.NewRightToUseEntity(
						utils.CreateJstTime(2020, 4, 1, 15, 11, 36, 123456),
						utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					),
				},
			)
			contractRep := db.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, conn)
			assert.NoError(t, err)
			// リロード
			savedContract, err := contractRep.GetById(savedContractId, conn)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedContract.RightToUses()[0], savedContract.BillingStartDate(), conn)
			assert.NoError(t, err)

			// 検証（日割り金額になる 31000 / 30 * 25）で端数切り捨て
			assert.Equal(t, "25833", billingAmount.String())
		})
	})

	t.Run("29日ある月", func(t *testing.T) {
		t.Run("使用権の有効期間が契約の課金単位期間_月払いとか年払いとか_の満了期間だと満額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成と課金開始日が同日）
			contractEntity := contract.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
				[]*contract.RightToUseEntity{
					contract.NewRightToUseEntity(
						utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
						utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
					),
				},
			)
			contractRep := db.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, conn)
			assert.NoError(t, err)
			// リロード
			savedContract, err := contractRep.GetById(savedContractId, conn)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedContract.RightToUses()[0], savedContract.BillingStartDate(), conn)
			assert.NoError(t, err)

			// 検証
			assert.Equal(t, "31000", billingAmount.String())
		})

		t.Run("使用権の期間が1ヶ月に足りず満了期間ではない場合日割り金額が返る", func(t *testing.T) {
			////// 準備
			// 契約を作成（契約作成から課金開始日まで5日間ある）
			contractEntity := contract.NewContractEntity(
				savedUserId,
				savedProductId,
				utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
				utils.CreateJstTime(2020, 2, 6, 0, 0, 0, 0),
				[]*contract.RightToUseEntity{
					contract.NewRightToUseEntity(
						utils.CreateJstTime(2020, 2, 1, 15, 11, 36, 123456),
						utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
					),
				},
			)
			contractRep := db.NewContractRepository()
			savedContractId, err := contractRep.Create(contractEntity, conn)
			assert.NoError(t, err)
			// リロード
			savedContract, err := contractRep.GetById(savedContractId, conn)

			// 実行
			billingAmount, err := billingDS.BillingAmount(savedContract.RightToUses()[0], savedContract.BillingStartDate(), conn)
			assert.NoError(t, err)

			// 検証（日割り金額になる 31000 / 29 * 24）で端数切り捨て
			assert.Equal(t, "25655", billingAmount.String())
		})
	})
}

func createProduct(price string) *product.ProductEntity {
	conn, err := db.GetConnection()
	if err != nil {
		panic("db接続エラー")
	}
	productEntity, err := product.NewProductEntity("商品", price)
	if err != nil {
		panic("データ作成エラー")
	}
	rep := db.NewProductRepository()
	id, err := rep.Save(productEntity, conn)
	if err != nil {
		panic("データ保存エラー")
	}
	entity, err := rep.GetById(id, conn)
	if err != nil {
		panic("データ取得エラー")
	}
	return entity
}
