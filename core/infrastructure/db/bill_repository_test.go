package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/bill"
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBillRepository_Create(t *testing.T) {
	t.Run("Bill集約を渡すと保存できる", func(t *testing.T) {
		////// 準備
		// 使用権作成
		rightToUse1Id, rightToUse2Id, userId := createRightToUseDataForTest()

		// 請求データ作成
		billEntity := bill.NewBillEntity(utils.CreateJstTime(2020, 8, 31, 0, 10, 0, 0), userId)
		err := billEntity.AddBillDetail(bill.NewBillDetailEntity(rightToUse1Id, decimal.NewFromInt(100)))
		assert.NoError(t, err)
		err = billEntity.AddBillDetail(bill.NewBillDetailEntity(rightToUse2Id, decimal.NewFromInt(1000)))
		assert.NoError(t, err)

		db, err := GetConnection()
		assert.NoError(t, err)
		defer db.Db.Close()

		rep := NewBillRepository()

		// 実行
		tran, err := db.Begin()
		assert.NoError(t, err)
		billId, err := rep.Create(billEntity, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証
		// billデータ取得
		var billMap BillMapper
		err = db.SelectOne(&billMap, "SELECT * FROM bills WHERE id = $1;", billId)
		assert.NoError(t, err)
		// billデータ検証
		assert.Equal(t, billId, billMap.Id)
		assert.True(t, billMap.BillingDate.Equal(utils.CreateJstTime(2020, 8, 31, 0, 10, 0, 0)))
		assert.False(t, billMap.PaymentConfirmedAt.Valid)
		assert.NotZero(t, billMap.CreatedAt)
		assert.NotZero(t, billMap.UpdatedAt)

		// billDetailデータ取得
		var billDetailsMap []BillDetailMapper
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
		rightToUse1Id, rightToUse2Id, userId := createRightToUseDataForTest()

		// 請求データ作成
		billEntity := bill.NewBillEntity(utils.CreateJstTime(2020, 8, 31, 0, 10, 0, 0), userId)
		err := billEntity.AddBillDetail(bill.NewBillDetailEntity(rightToUse1Id, decimal.NewFromInt(100)))
		assert.NoError(t, err)
		err = billEntity.AddBillDetail(bill.NewBillDetailEntity(rightToUse2Id, decimal.NewFromInt(1000)))
		assert.NoError(t, err)

		db, err := GetConnection()
		assert.NoError(t, err)
		defer db.Db.Close()

		// 請求データ保存
		rep := NewBillRepository()
		tran, err := db.Begin()
		assert.NoError(t, err)
		billId, err := rep.Create(billEntity, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		// データ取得
		actual, err := rep.GetById(billId, db)
		assert.NoError(t, err)

		// 検証
		assert.Equal(t, billId, actual.Id())
		assertBillAgg(t, billEntity, actual)
	})
}

func TestBillRepository_GetByUserId(t *testing.T) {
	t.Run("UserIdを渡すとデータを取得できる", func(t *testing.T) {
		t.Run("billが複数あれば複数がスライスで取れる", func(t *testing.T) {
			////// 準備
			// 使用権作成
			rightToUse1Id, rightToUse2Id, userId := createRightToUseDataForTest()

			// 請求データ作成
			billEntity1 := bill.NewBillEntity(utils.CreateJstTime(2020, 8, 1, 0, 10, 0, 0), userId)
			err := billEntity1.AddBillDetail(bill.NewBillDetailEntity(rightToUse1Id, decimal.NewFromInt(100)))
			assert.NoError(t, err)

			billEntity2 := bill.NewBillEntity(utils.CreateJstTime(2020, 9, 1, 0, 10, 0, 0), userId)
			err = billEntity2.AddBillDetail(bill.NewBillDetailEntity(rightToUse2Id, decimal.NewFromInt(1000)))
			assert.NoError(t, err)

			db, err := GetConnection()
			assert.NoError(t, err)
			defer db.Db.Close()

			// 請求データ保存
			rep := NewBillRepository()
			tran, err := db.Begin()
			assert.NoError(t, err)
			billId1, err := rep.Create(billEntity1, tran)
			assert.NoError(t, err)
			billId2, err := rep.Create(billEntity2, tran)
			assert.NoError(t, err)
			err = tran.Commit()
			assert.NoError(t, err)

			// データ取得
			actual, err := rep.GetByUserId(userId, db)
			assert.NoError(t, err)

			// 検証
			assert.Len(t, actual, 2) // userIdのbillは2つあるので。
			// Idを検証
			assert.Equal(t, billId1, actual[0].Id())
			assert.Equal(t, billId2, actual[1].Id())
			// その他の要素の検証
			expect := []*bill.BillEntity{
				billEntity1,
				billEntity2,
			}
			for i := range actual {
				assertBillAgg(t, expect[i], actual[i])
			}
		})

		t.Run("一つのbillに複数のbillDetailがある場合そのようにデータが返却される", func(t *testing.T) {
			////// 準備
			// 使用権作成
			rightToUse1Id, rightToUse2Id, userId := createRightToUseDataForTest()

			// 請求データ作成
			billEntity1 := bill.NewBillEntity(utils.CreateJstTime(2020, 8, 1, 0, 10, 0, 0), userId)
			err := billEntity1.AddBillDetail(bill.NewBillDetailEntity(rightToUse1Id, decimal.NewFromInt(100)))
			err = billEntity1.AddBillDetail(bill.NewBillDetailEntity(rightToUse2Id, decimal.NewFromInt(1000)))
			assert.NoError(t, err)

			db, err := GetConnection()
			assert.NoError(t, err)
			defer db.Db.Close()

			// 請求データ保存
			rep := NewBillRepository()
			tran, err := db.Begin()
			assert.NoError(t, err)
			billId1, err := rep.Create(billEntity1, tran)
			assert.NoError(t, err)
			err = tran.Commit()
			assert.NoError(t, err)

			// データ取得
			actual, err := rep.GetByUserId(userId, db)
			assert.NoError(t, err)

			// 検証
			assert.Len(t, actual, 1) // userIdのbillは1つあるので。
			// Idを検証
			assert.Equal(t, billId1, actual[0].Id())
			// その他の要素の検証
			expect := []*bill.BillEntity{
				billEntity1,
			}
			for i := range actual {
				assertBillAgg(t, expect[i], actual[i])
			}
		})
	})
}

// 請求集約のアサーション。IdやCreatedAtやUpdatedAtなどはテストしにくいためしてない
func assertBillAgg(t *testing.T, expect, actual *bill.BillEntity) {
	assert.True(t, expect.BillingDate().Equal(actual.BillingDate()))
	assert.Equal(t, expect.UserId(), actual.UserId())
	expectConfirmedAt, isNull, err := expect.PaymentConfirmedAt()
	assert.NoError(t, err)
	assert.True(t, isNull)
	assert.Zero(t, expectConfirmedAt)
	actualConfirmedAt, isNull, err := actual.PaymentConfirmedAt()
	assert.NoError(t, err)
	assert.True(t, isNull)
	assert.Zero(t, actualConfirmedAt)
	assert.NotZero(t, actual.CreatedAt())
	assert.NotZero(t, actual.UpdatedAt())

	expectDetails := expect.BillDetails()
	actualDetails := actual.BillDetails()
	assert.Equal(t, len(expectDetails), len(actualDetails))
	for i := range expectDetails {
		assert.NotZero(t, actualDetails[i].Id())
		assert.Equal(t, expectDetails[i].RightToUseId(), actualDetails[i].RightToUseId())
		expectBillingAmount := expectDetails[i].BillingAmount()
		actualBillingAmount := actualDetails[i].BillingAmount()
		assert.Equal(t, expectBillingAmount.String(), actualBillingAmount.String())
		assert.NotZero(t, actualDetails[i].CreatedAt())
		assert.NotZero(t, actualDetails[i].UpdatedAt())
	}
}

func createRightToUseDataForTest() (rightToUse1Id, rightToUse2Id, userId int) {
	db, err := GetConnection()
	if err != nil {
		panic("dbコネクション取得失敗")
	}

	// 商品データ作成
	productEntity1, err := product.NewProductEntity("商品", "100")
	if err != nil {
		panic("商品データ作成失敗")
	}
	productEntity2, err := product.NewProductEntity("商品", "1000")
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
	userEntity, err := user.NewUserIndividualEntity("個人請求太郎")
	if err != nil {
		panic("ユーザーデータ作成失敗")
	}
	userRep := NewUserRepository()
	userId, err = userRep.SaveUserIndividual(userEntity, db)
	if err != nil {
		panic("ユーザーデータ登録失敗")
	}

	// 使用権データ作成
	rightToUseEntity1 := contract.NewRightToUseEntity(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0))
	rightToUseEntity2 := contract.NewRightToUseEntity(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0))
	// 契約データ作成
	contractEntity1 := contract.NewContractEntity(
		userId,
		product1Id,
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
		[]*contract.RightToUseEntity{
			rightToUseEntity1,
		},
	)
	contractEntity2 := contract.NewContractEntity(
		userId,
		product2Id,
		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
		[]*contract.RightToUseEntity{
			rightToUseEntity2,
		},
	)
	contractRep := NewContractRepository()
	contract1Id, err := contractRep.Create(contractEntity1, db)
	if err != nil {
		panic("契約データ登録失敗")
	}
	contract2Id, err := contractRep.Create(contractEntity2, db)
	if err != nil {
		panic("契約データ登録失敗")
	}
	// 登録データ再読込
	reloadedContract1, err := contractRep.GetById(contract1Id, db)
	if err != nil {
		panic("契約データ再読込失敗")
	}
	reloadedContract2, err := contractRep.GetById(contract2Id, db)
	if err != nil {
		panic("契約データ再読込失敗")
	}
	return reloadedContract1.RightToUses()[0].Id(), reloadedContract2.RightToUses()[0].Id(), userId
}
