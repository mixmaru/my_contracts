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

// トランザクションが正しく動作しているかテスト
//func TestUserRepository_Transaction(t *testing.T) {
//	// db接続
//	dbMap, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer dbMap.Db.Close()
//
//	t.Run("コミットするとデータ保存されている", func(t *testing.T) {
//		// transaction開始
//		tran, err := dbMap.Begin()
//		assert.NoError(t, err)
//
//		//データ保存
//		user, err := entities.NewUserIndividualEntity("個人太郎")
//		assert.NoError(t, err)
//		repo := NewUserRepository()
//		savedId, err := repo.SaveUserIndividual(user, tran)
//		assert.NoError(t, err)
//
//		// コミット
//		err = tran.Commit()
//		assert.NoError(t, err)
//
//		// データ取得できる
//		_, err = repo.GetUserIndividualById(savedId, dbMap)
//		assert.NoError(t, err) // sql: no rows in result set エラーが起こらなければ、データが保存されている
//	})
//
//	t.Run("ロールバックするとデータ保存されていない", func(t *testing.T) {
//		// transaction開始
//		tran, err := dbMap.Begin()
//		assert.NoError(t, err)
//
//		//データ保存
//		user, err := entities.NewUserIndividualEntity("個人太郎")
//		assert.NoError(t, err)
//		repo := NewUserRepository()
//		savedId, err := repo.SaveUserIndividual(user, tran)
//		assert.NoError(t, err)
//
//		// ロールバック
//		err = tran.Rollback()
//		assert.NoError(t, err)
//
//		// データ取得できない
//		user, err = repo.GetUserIndividualById(savedId, dbMap)
//		assert.Nil(t, user)
//	})
//}
//
//func TestUserRepository_SaveUserIndividual(t *testing.T) {
//	// 登録用データ作成
//	user, err := entities.NewUserIndividualEntity("個人太郎")
//	assert.NoError(t, err)
//
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//
//	// 実行
//	repo := NewUserRepository()
//	_, err = repo.SaveUserIndividual(user, db)
//	assert.NoError(t, err)
//}
//
//func TestUserRepository_GetUserIndividualById(t *testing.T) {
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//
//	//　事前にデータ登録する
//	user, err := entities.NewUserIndividualEntity("個人太郎")
//	assert.NoError(t, err)
//	repo := NewUserRepository()
//	savedId, err := repo.SaveUserIndividual(user, db)
//	assert.NoError(t, err)
//
//	// idで取得して検証
//	t.Run("データがある時_idでデータが取得できる", func(t *testing.T) {
//		result, err := repo.GetUserIndividualById(savedId, db)
//		assert.NoError(t, err)
//		assert.Equal(t, user.Name(), result.Name())
//		assert.NotEqual(t, time.Time{}, result.CreatedAt())
//		assert.NotEqual(t, time.Time{}, result.UpdatedAt())
//	})
//
//	t.Run("データが無い時_nilが返る", func(t *testing.T) {
//		user, err := repo.GetUserIndividualById(-1, db)
//		assert.NoError(t, err)
//		assert.Nil(t, user)
//	})
//}
//
//func TestUserRepository_GetUserCorporationById(t *testing.T) {
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//
//	//　事前にデータ登録する
//	savingUser, err := entities.NewUserCorporationEntity("イケてる会社", "担当　太郎", "社長　太郎")
//	assert.NoError(t, err)
//
//	repo := NewUserRepository()
//	savedId, err := repo.SaveUserCorporation(savingUser, db)
//	assert.NoError(t, err)
//
//	// idで取得して検証
//	t.Run("データがある時_idでデータが取得できる", func(t *testing.T) {
//		result, err := repo.GetUserCorporationById(savedId, db)
//		assert.NoError(t, err)
//		assert.Equal(t, savedId, result.Id())
//		assert.Equal(t, "イケてる会社", result.CorporationName())
//		assert.Equal(t, "担当　太郎", result.ContactPersonName())
//		assert.Equal(t, "社長　太郎", result.PresidentName())
//		assert.NotEqual(t, time.Time{}, result.CreatedAt())
//		assert.NotEqual(t, time.Time{}, result.UpdatedAt())
//	})
//
//	t.Run("データが無い時_nilが返る", func(t *testing.T) {
//		result, err := repo.GetUserCorporationById(-1, db)
//		assert.NoError(t, err)
//		assert.Nil(t, result)
//	})
//}
//
//func TestUserRepository_SaveUserCorporation(t *testing.T) {
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//
//	// 保存するデータ作成
//	user, err := entities.NewUserCorporationEntity("イケてる会社", "担当太郎", "社長次郎")
//	assert.NoError(t, err)
//
//	// 保存実行
//	repo := NewUserRepository()
//	_, err = repo.SaveUserCorporation(user, db)
//	assert.NoError(t, err)
//}
//
//func TestUserRepository_getUserCorporationViewById(t *testing.T) {
//	// db接続
//	dbMap, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer dbMap.Db.Close()
//
//	// 事前にデータ登録
//	user, err := entities.NewUserCorporationEntity("イケてる会社", "担当太郎", "社長次郎")
//	assert.NoError(t, err)
//	repo := NewUserRepository()
//	savedId, err := repo.SaveUserCorporation(user, dbMap)
//	assert.NoError(t, err)
//
//	// idで取得する
//	result, err := repo.getUserCorporationEntityById(savedId, &entities.UserCorporationEntity{}, dbMap)
//	assert.NoError(t, err)
//
//	// 検証
//	assert.Equal(t, result.Id(), savedId)
//	assert.Equal(t, "イケてる会社", result.CorporationName())
//	assert.Equal(t, "担当太郎", result.ContactPersonName())
//	assert.Equal(t, "社長次郎", result.PresidentName())
//	assert.NotEqual(t, time.Time{}, result.CreatedAt())
//	assert.NotEqual(t, time.Time{}, result.UpdatedAt())
//}
//
//func TestUserRepository_GetUserById(t *testing.T) {
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//
//	//　事前にデータ登録する。個人顧客
//	userIndividual, err := entities.NewUserIndividualEntity("個人太郎")
//	assert.NoError(t, err)
//	repo := NewUserRepository()
//	savedIndividualId, err := repo.SaveUserIndividual(userIndividual, db)
//	assert.NoError(t, err)
//
//	//　事前にデータ登録する。法人顧客
//	userCorporation, err := entities.NewUserCorporationEntity("イケてる会社", "担当太郎", "社長太郎")
//	assert.NoError(t, err)
//	repo = NewUserRepository()
//	savedCorporationId, err := repo.SaveUserCorporation(userCorporation, db)
//	assert.NoError(t, err)
//
//	// idで取得して検証
//	t.Run("個人顧客データ取得", func(t *testing.T) {
//		t.Run("データがある時_interface{}型でUserIndividualEntityが返る", func(t *testing.T) {
//			result, err := repo.GetUserById(savedIndividualId, db)
//			assert.NoError(t, err)
//
//			loadedIndividual, ok := result.(*entities.UserIndividualEntity)
//			assert.True(t, ok)
//
//			assert.Equal(t, savedIndividualId, loadedIndividual.Id())
//			assert.Equal(t, "個人太郎", loadedIndividual.Name())
//			assert.NotZero(t, loadedIndividual.CreatedAt())
//			assert.NotZero(t, loadedIndividual.UpdatedAt())
//		})
//
//		t.Run("データが無い時_nilが返る", func(t *testing.T) {
//			user, err := repo.GetUserById(-1, db)
//			assert.NoError(t, err)
//			assert.Nil(t, user)
//		})
//	})
//
//	t.Run("法人顧客データ取得", func(t *testing.T) {
//		t.Run("データがある時_interface{}型でUserCorporationEntityが返る", func(t *testing.T) {
//			result, err := repo.GetUserById(savedCorporationId, db)
//			assert.NoError(t, err)
//
//			loadedCorporation, ok := result.(*entities.UserCorporationEntity)
//			assert.True(t, ok)
//
//			assert.Equal(t, savedCorporationId, loadedCorporation.Id())
//			assert.Equal(t, "イケてる会社", loadedCorporation.CorporationName())
//			assert.Equal(t, "担当太郎", loadedCorporation.ContactPersonName())
//			assert.Equal(t, "社長太郎", loadedCorporation.PresidentName())
//			assert.NotZero(t, loadedCorporation.CreatedAt())
//			assert.NotZero(t, loadedCorporation.UpdatedAt())
//		})
//
//		t.Run("データが無い時_nilが返る", func(t *testing.T) {
//			user, err := repo.GetUserById(-1, db)
//			assert.NoError(t, err)
//			assert.Nil(t, user)
//		})
//	})
//}

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
