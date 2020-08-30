package repositories

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBillRepository_Create(t *testing.T) {
	t.Run("Bill集約を渡すと保存できる", func(t *testing.T) {
		////// 準備
		// 使用権作成
		rightToUse1Id, rightToUse2Id := createRightToUseDataForTest()

		// 請求データ作成
		billAgg := entities.NewBillingAggregation(utils.CreateJstTime(2020, 8, 31, 0, 10, 0, 0))
		err := billAgg.AddBillDetail(entities.NewBillingDetailEntity(1, rightToUse1Id, decimal.NewFromInt(100)))
		assert.NoError(t, err)
		err = billAgg.AddBillDetail(entities.NewBillingDetailEntity(2, rightToUse2Id, decimal.NewFromInt(1000)))
		assert.NoError(t, err)

		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		defer db.Db.Close()

		rep := NewBillRepository()

		// 実行
		tran, err := db.Begin()
		assert.NoError(t, err)
		billId, err := rep.Create(billAgg, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証
		// billデータ取得
		var billMap data_mappers.BillMapper
		err = db.SelectOne(&billMap, "SELECT * FROM bills WHERE id = $1;", billId)
		assert.NoError(t, err)
		// billデータ検証
		assert.Equal(t, billId, billMap.Id)
		assert.True(t, billMap.BillingDate.Equal(utils.CreateJstTime(2020, 8, 31, 0, 10, 0, 0)))
		assert.False(t, billMap.PaymentConfirmedAt.Valid)
		assert.NotZero(t, billMap.CreatedAt)
		assert.NotZero(t, billMap.UpdatedAt)

		// billDetailデータ取得
		var billDetailsMap []data_mappers.BillDetailMapper
		_, err = db.Select(&billDetailsMap, "SELECT * FROM bill_details WHERE bill_id = $1 ORDER BY order_num", billId)
		assert.NoError(t, err)
		// billDetailデータ検証
		assert.Len(t, billDetailsMap, 2)

		assert.NotZero(t, billDetailsMap[0].Id)
		assert.Equal(t, billId, billDetailsMap[0].BillId)
		assert.Equal(t, 1, billDetailsMap[0].OrderNum)
		assert.Equal(t, rightToUse1Id, billDetailsMap[0].RightToUseId)
		assert.Equal(t, "100", billDetailsMap[0].BillingAmount.String())
		assert.NotZero(t, billDetailsMap[0].CreatedAt)
		assert.NotZero(t, billDetailsMap[0].UpdatedAt)

		assert.NotZero(t, billDetailsMap[1].Id)
		assert.Equal(t, billId, billDetailsMap[1].BillId)
		assert.Equal(t, 2, billDetailsMap[1].OrderNum)
		assert.Equal(t, rightToUse2Id, billDetailsMap[1].RightToUseId)
		assert.Equal(t, "1000", billDetailsMap[1].BillingAmount.String())
		assert.NotZero(t, billDetailsMap[1].CreatedAt)
		assert.NotZero(t, billDetailsMap[1].UpdatedAt)
	})
}

func TestBillRepository_GetById(t *testing.T) {
	t.Run("idを渡すとデータを取得できる", func(t *testing.T) {
		////// 準備
		// 使用権作成
		rightToUse1Id, rightToUse2Id := createRightToUseDataForTest()

		// 請求データ作成
		billAgg := entities.NewBillingAggregation(utils.CreateJstTime(2020, 8, 31, 0, 10, 0, 0))
		err := billAgg.AddBillDetail(entities.NewBillingDetailEntity(1, rightToUse1Id, decimal.NewFromInt(100)))
		assert.NoError(t, err)
		err = billAgg.AddBillDetail(entities.NewBillingDetailEntity(2, rightToUse2Id, decimal.NewFromInt(1000)))
		assert.NoError(t, err)

		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		defer db.Db.Close()

		// 請求データ保存
		rep := NewBillRepository()
		tran, err := db.Begin()
		assert.NoError(t, err)
		billId, err := rep.Create(billAgg, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		// データ取得
		actual, err := rep.GetById(billId, db)
		assert.NoError(t, err)

		assert.Equal(t, billId, actual.Id())
		assert.True(t, actual.BillingDate().Equal(utils.CreateJstTime(2020, 8, 31, 0, 10, 0, 0)))
		confirmedAt, isNull, err := actual.PaymentConfirmedAt()
		assert.NoError(t, err)
		assert.True(t, isNull)
		assert.Zero(t, confirmedAt)
		assert.NotZero(t, actual.CreatedAt())
		assert.NotZero(t, actual.UpdatedAt())

		details := actual.BillDetails()
		assert.Len(t, details, 2)

		assert.NotZero(t, details[0].Id())
		assert.Equal(t, 1, details[0].OrderNum())
		assert.Equal(t, rightToUse1Id, details[0].RightToUseId())
		billingAmount1 := details[0].BillingAmount()
		assert.Equal(t, "100", billingAmount1.String())
		assert.NotZero(t, details[0].CreatedAt())
		assert.NotZero(t, details[0].UpdatedAt())

		assert.NotZero(t, details[1].Id())
		assert.Equal(t, 2, details[1].OrderNum())
		assert.Equal(t, rightToUse2Id, details[1].RightToUseId())
		billingAmount2 := details[1].BillingAmount()
		assert.Equal(t, "1000", billingAmount2.String())
		assert.NotZero(t, details[1].CreatedAt())
		assert.NotZero(t, details[1].UpdatedAt())
	})
}

func createRightToUseDataForTest() (rightToUse1Id, rightToUse2Id int) {
	db, err := db_connection.GetConnection()
	if err != nil {
		panic("dbコネクション取得失敗")
	}

	// 商品データ作成
	productEntity1, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), "100")
	if err != nil {
		panic("商品データ作成失敗")
	}
	productEntity2, err := entities.NewProductEntity(utils.CreateUniqProductNameForTest(), "1000")
	if err != nil {
		panic("商品データ作成失敗")
	}
	productRep := NewProductRepository()
	product1Id, err := productRep.Save(productEntity1, db)
	if err != nil {
		panic("商品データ登録失敗")
	}
	product2Id, err := productRep.Save(productEntity2, db)
	if err != nil {
		panic("商品データ登録失敗")
	}

	// userデータ作成
	userEntity, err := entities.NewUserIndividualEntity("個人請求太郎")
	if err != nil {
		panic("ユーザーデータ作成失敗")
	}
	userRep := NewUserRepository()
	userId, err := userRep.SaveUserIndividual(userEntity, db)
	if err != nil {
		panic("ユーザーデータ登録失敗")
	}

	// 契約データ作成
	contractEntity1 := entities.NewContractEntity(userId, product1Id, utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0))
	contractEntity2 := entities.NewContractEntity(userId, product2Id, utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0))
	contractRep := NewContractRepository()
	contract1Id, err := contractRep.Create(contractEntity1, db)
	if err != nil {
		panic("契約データ登録失敗")
	}
	contract2Id, err := contractRep.Create(contractEntity2, db)
	if err != nil {
		panic("契約データ登録失敗")
	}

	// 使用権データ作成
	rightToUseEntity1 := entities.NewRightToUseEntity(contract1Id, utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0))
	rightToUseEntity2 := entities.NewRightToUseEntity(contract2Id, utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0))
	rightToUseRep := NewRightToUseRepository()
	rightToUse1Id, err = rightToUseRep.Create(rightToUseEntity1, db)
	if err != nil {
		panic("使用権データ登録失敗")
	}
	rightToUse2Id, err = rightToUseRep.Create(rightToUseEntity2, db)
	if err != nil {
		panic("使用権データ登録失敗")
	}
	return rightToUse1Id, rightToUse2Id
}
