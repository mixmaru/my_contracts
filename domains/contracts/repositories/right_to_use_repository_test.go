package repositories

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestRightToUseRepository_Create(t *testing.T) {
//	r := NewRightToUseRepository()
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//
//	////// テスト用契約を作成する
//	// 契約の作成
//	_ = createPreparedContractData(
//		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
//		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
//		db,
//	)
//
//	t.Run("権利エンティティとdbコネクションを渡すとDBへ新規保存されて_保存Idを返す", func(t *testing.T) {
//		// 準備
//		entity := entities.NewRightToUseEntity(
//			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
//			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
//		)
//
//		// 実行
//		savedId, err := r.Create(entity, db)
//		assert.NoError(t, err)
//
//		//検証
//		assert.NotZero(t, savedId)
//	})
//}
//
//func TestRightToUseRepository_GetById(t *testing.T) {
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer db.Db.Close()
//
//	// 事前に使用権を登録する
//	r := NewRightToUseRepository()
//	_ = createPreparedContractData(
//		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
//		utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0),
//		db,
//	)
//	rightToUseEntity := entities.NewRightToUseEntity(
//		utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
//		utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
//	)
//
//	// 実行
//	savedId, err := r.Create(rightToUseEntity, db)
//	assert.NoError(t, err)
//
//	t.Run("データがあればIdを渡すとデータが取得できる", func(t *testing.T) {
//		// 実行
//		actual, err := r.GetById(savedId, db)
//		assert.NoError(t, err)
//
//		// 検証
//		assert.Equal(t, savedId, actual.Id())
//		//assert.Equal(t, savedContractId, actual.ContractId())
//		assert.True(t, actual.ValidFrom().Equal(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)))
//		assert.True(t, actual.ValidTo().Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
//	})
//
//	t.Run("データがなければidを渡すとnilが返る", func(t *testing.T) {
//		// 実行
//		actual, err := r.GetById(-100, db)
//		assert.NoError(t, err)
//
//		// 検証
//		assert.Nil(t, actual)
//	})
//}

//func TestRightToUseRepository_GetBillingTargetByBillingDate(t *testing.T) {
//	db, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//
//	r := NewRightToUseRepository()
//	t.Run("請求実行日を渡すとその日以前で請求実行をしていない使用権データがuserId、rightToUserId順で全て返る", func(t *testing.T) {
//		// 準備
//		// 事前に対象になる使用権を削除しておく
//		tran, err := db.Begin()
//		assert.NoError(t, err)
//
//		deleteTestRightToUseData(tran)
//
//		// 2userに対して、契約の課金開始日が6/11の以下の使用権を作成する。
//		// 6/1 ~ 6/30 の使用権（未請求)
//		// 6/1 ~ 6/30 の使用権（請求済）
//		// 7/1 ~ 7/31 の使用権（未請求)
//		// 8/1 ~ 8/31 の使用権（未請求）
//		rightToUseIds1, _ := createBillTestData(tran)
//		rightToUseIds2, _ := createBillTestData(tran)
//
//		// 実行
//		billingDate := utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)
//		actual, err := r.GetBillingTargetByBillingDate(billingDate, tran)
//		assert.NoError(t, err)
//
//		err = tran.Commit()
//		assert.NoError(t, err)
//
//		// 検証
//		// 6/1 ~ 6/30 の使用権（未請求) => 取得される（未請求だから）
//		// 6/1 ~ 6/30 の使用権（請求済）=> 取得されない（請求実行済だから）
//		// 7/1 ~ 7/31 の使用権（未請求) => 取得される（未請求だから）
//		// 8/1 ~ 8/31 の使用権（未請求） => 取得されない（まだ請求しない使用権だから）
//		assert.Len(t, actual, 4)
//		assert.Equal(t, rightToUseIds1[0], actual[0].Id())
//		assert.Equal(t, rightToUseIds1[2], actual[1].Id())
//		assert.Equal(t, rightToUseIds2[0], actual[2].Id())
//		assert.Equal(t, rightToUseIds2[2], actual[3].Id())
//	})
//
//	t.Run("渡した請求実行日が契約の課金開始日以前である使用権は返却データに含まれない_課金開始日以前の使用権には請求が発生しない", func(t *testing.T) {
//		tran, err := db.Begin()
//		assert.NoError(t, err)
//
//		// 事前に影響のあるデータを削除
//		deleteTestRightToUseData(tran)
//
//		// 準備（以下の使用権データを作成する）
//		// 6/1 ~ 6/30 の使用権（未請求）=> 取得されない（契約の課金開始日が6/11だから）
//		// 6/1 ~ 6/30 の使用権（請求済）=> 取得されない
//		// 7/1 ~ 7/31 の使用権（未請求）=> 取得されない
//		// 8/1 ~ 8/31 の使用権（未請求）=> 取得されない
//		_, _ = createBillTestData(tran)
//
//		// 実行
//		billingDate := utils.CreateJstTime(2020, 6, 10, 0, 0, 0, 0)
//		actual, err := r.GetBillingTargetByBillingDate(billingDate, tran)
//		assert.NoError(t, err)
//
//		err = tran.Commit()
//		assert.NoError(t, err)
//
//		// 検証
//		assert.Len(t, actual, 0)
//	})
//
//	t.Run("渡した請求実行日が契約の課金開始日ちょうどである使用権は返却データに含まれる", func(t *testing.T) {
//		tran, err := db.Begin()
//		assert.NoError(t, err)
//
//		// 事前に影響のあるデータを削除
//		deleteTestRightToUseData(tran)
//
//		// 準備（以下の使用権データを作成する）
//		// 6/1 ~ 6/30 の使用権（未請求）=> 取得される（契約の課金開始日が6/11だから）
//		// 6/1 ~ 6/30 の使用権（請求済）=> 取得されない
//		// 7/1 ~ 7/31 の使用権（未請求）=> 取得されない
//		// 8/1 ~ 8/31 の使用権（未請求）=> 取得されない
//		rightToUseIds, _ := createBillTestData(tran)
//
//		// 実行
//		billingDate := utils.CreateJstTime(2020, 6, 11, 0, 0, 0, 0)
//		actual, err := r.GetBillingTargetByBillingDate(billingDate, tran)
//		assert.NoError(t, err)
//
//		err = tran.Commit()
//		assert.NoError(t, err)
//
//		// 検証
//		assert.Len(t, actual, 1)
//		assert.Equal(t, rightToUseIds[0], actual[0].Id())
//	})
//}

func TestRightToUseRepository_GetRecurTargets(t *testing.T) {
	t.Run("2020/6/1を渡すと使用期間終了日が6/1 ~ 6/5でかつ次の期間の使用権がまだない使用権が返る", func(t *testing.T) {
		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		defer db.Db.Close()
		tran, err := db.Begin()
		assert.NoError(t, err)

		////// 準備
		// 事前に影響するデータを削除しておく
		_, err = tran.Exec(
			"DELETE FROM bill_details",
		)
		assert.NoError(t, err)
		_, err = tran.Exec(
			"DELETE FROM right_to_use_active",
		)
		assert.NoError(t, err)
		_, err = tran.Exec(
			"DELETE FROM right_to_use",
		)
		assert.NoError(t, err)
		// テスト用データの登録
		_ = createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*entities.RightToUseEntity{
				entities.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 5, 31, 0, 0, 0, 0),
				),
			},
			tran,
		)
		contractB := createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*entities.RightToUseEntity{
				entities.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
				),
			},
			tran,
		)
		contractC := createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*entities.RightToUseEntity{
				entities.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 6, 5, 0, 0, 0, 0),
				),
			},
			tran,
		)
		_ = createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*entities.RightToUseEntity{
				entities.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 6, 5, 0, 0, 0, 0),
				),
				entities.NewRightToUseEntity(
					utils.CreateJstTime(2020, 6, 6, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 7, 5, 0, 0, 0, 0),
				),
			},
			tran,
		)
		_ = createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*entities.RightToUseEntity{
				entities.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 6, 6, 0, 0, 0, 0),
				),
			},
			tran,
		)
		// 使用権を作成（終了日が5/31, 6/1, 6/5, 6/5（ただし次の使用権データがある） 6/6）
		rightToUseRep := NewRightToUseRepository()
		//rightToUse1 := entities.NewRightToUseEntity(
		//	utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
		//	utils.CreateJstTime(2020, 5, 31, 0, 0, 0, 0),
		//)
		//_, err = rightToUseRep.Create(rightToUse1, tran)
		//assert.NoError(t, err)
		//rightToUse2 := entities.NewRightToUseEntity(
		//	utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
		//	utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
		//)
		//rightToUse2Id, err := rightToUseRep.Create(rightToUse2, tran)
		//assert.NoError(t, err)
		//rightToUse3 := entities.NewRightToUseEntity(
		//	utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
		//	utils.CreateJstTime(2020, 6, 5, 0, 0, 0, 0),
		//)
		//rightToUse3Id, err := rightToUseRep.Create(rightToUse3, tran)
		//assert.NoError(t, err)
		//rightToUse4 := entities.NewRightToUseEntity(
		//	utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
		//	utils.CreateJstTime(2020, 6, 5, 0, 0, 0, 0),
		//)
		//_, err = rightToUseRep.Create(rightToUse4, tran)
		//assert.NoError(t, err)
		//rightToUse4Next := entities.NewRightToUseEntity(
		//	//contractIdD,
		//	utils.CreateJstTime(2020, 6, 6, 0, 0, 0, 0),
		//	utils.CreateJstTime(2020, 7, 5, 0, 0, 0, 0),
		//)
		//_, err = rightToUseRep.Create(rightToUse4Next, tran)
		//assert.NoError(t, err)
		//rightToUse5 := entities.NewRightToUseEntity(
		//	//contractIdE,
		//	utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
		//	utils.CreateJstTime(2020, 6, 6, 0, 0, 0, 0),
		//)
		//_, err = rightToUseRep.Create(rightToUse5, tran)
		//assert.NoError(t, err)

		////// 実行
		expects, err := rightToUseRep.GetRecurTargets(utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0), tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, expects, 2)
		assert.Equal(t, contractB.RightToUses()[0].Id(), expects[0].Id())
		assert.Equal(t, contractC.RightToUses()[0].Id(), expects[1].Id())
	})
}
