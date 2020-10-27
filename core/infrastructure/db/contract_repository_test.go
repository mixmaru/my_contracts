package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

func TestContractRepository_Create(t *testing.T) {
	db, err := GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// userを作成
	userRepository := NewUserRepository()
	userEntity, err := user.NewUserIndividualEntity("担当太郎")
	assert.NoError(t, err)
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, db)
	assert.NoError(t, err)

	// 商品を登録
	productRepository := NewProductRepository()
	productEntity, err := product.NewProductEntity("商品", "1000")
	assert.NoError(t, err)
	savedProductId, err := productRepository.Save(productEntity, db)
	assert.NoError(t, err)

	t.Run("UserIdとProductIdと契約日と課金開始日を渡すと契約が新規作成される", func(t *testing.T) {
		////// 準備
		// 使用権データ作成
		rightToUses := make([]*contract.RightToUseEntity, 0, 2)
		rightToUse1 := contract.NewRightToUseEntity(
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
		)
		rightToUse2 := contract.NewRightToUseEntity(
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
		)
		rightToUses = append(rightToUses, rightToUse1, rightToUse2)
		// 契約データ作成
		contractDate := utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)
		billingStartDate := utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0)
		contractEntity := contract.NewContractEntity(savedUserId, savedProductId, contractDate, billingStartDate, rightToUses)

		////// 実行
		contractRepository := NewContractRepository()
		savedContractId, err := contractRepository.Create(contractEntity, db)

		////// 検証
		assert.NoError(t, err)
		assert.NotZero(t, savedContractId)
		// 使用権が作られているかチェック
		count, err := db.SelectInt(`
SELECT COUNT(1)
FROM right_to_use_active rtua
    INNER JOIN right_to_use rtu ON rtua.right_to_use_id = rtu.id
    INNER JOIN contracts c ON c.id = rtu.contract_id
WHERE c.id = $1
GROUP BY contract_id
;`, savedContractId)
		assert.NoError(t, err)
		assert.Equal(t, 2, int(count))

	})

	t.Run("存在しないuserIdで作成されようとしたとき_エラーが出る", func(t *testing.T) {
		contractRepository := NewContractRepository()
		contractDate := utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)
		billingStartDate := utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0)
		contractEntity := contract.NewContractEntity(-100, savedProductId, contractDate, billingStartDate, []*contract.RightToUseEntity{})

		savedContractId, err := contractRepository.Create(contractEntity, db)

		assert.Error(t, err)
		assert.Zero(t, savedContractId)
	})

	t.Run("存在しないproductIDで作成されようとしたとき_エラーが出る", func(t *testing.T) {
		contractRepository := NewContractRepository()
		contractDate := utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)
		billingStartDate := utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0)
		contractEntity := contract.NewContractEntity(savedUserId, -100, contractDate, billingStartDate, []*contract.RightToUseEntity{})

		savedContractId, err := contractRepository.Create(contractEntity, db)

		assert.Error(t, err)
		assert.Zero(t, savedContractId)
	})
}

func TestContractRepository_GetById(t *testing.T) {
	db, err := GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// userを作成
	userRepository := NewUserRepository()
	userEntity, err := user.NewUserIndividualEntity("担当太郎")
	assert.NoError(t, err)
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, db)
	assert.NoError(t, err)

	// 商品を登録
	productRepository := NewProductRepository()
	productEntity, err := product.NewProductEntity("商品", "1000")
	assert.NoError(t, err)
	savedProductId, err := productRepository.Save(productEntity, db)
	assert.NoError(t, err)

	t.Run("データがある時_Idで契約エンティティを返す", func(t *testing.T) {
		r := NewContractRepository()
		// データ登録
		contractEntity := contract.NewContractEntity(
			savedUserId,
			savedProductId,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
				),
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0),
				),
			},
		)
		savedId, err := r.Create(contractEntity, db)
		assert.NoError(t, err)

		// データ取得（使用権idを取得するため）
		loadedContract, err := r.GetById(savedId, db)
		assert.NoError(t, err)
		// 1つめの使用権は請求済にしておく
		bill := entities.NewBillingAggregation(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0), savedUserId)
		err = bill.AddBillDetail(entities.NewBillingDetailEntity(loadedContract.RightToUses()[0].Id(), decimal.NewFromInt(10000)))
		assert.NoError(t, err)
		billRep := repositories.NewBillRepository()
		_, err = billRep.Create(bill, db)
		assert.NoError(t, err)
		// データ取得（請求済データを反映するため）
		loadedContract, err = r.GetById(savedId, db)
		assert.NoError(t, err)

		// contractテスト
		assert.Equal(t, savedId, loadedContract.Id())
		assert.Equal(t, savedUserId, loadedContract.UserId())
		assert.Equal(t, savedProductId, loadedContract.ProductId())
		assert.NotZero(t, loadedContract.CreatedAt())
		assert.NotZero(t, loadedContract.UpdatedAt())
		// rightToUse
		rightToUses := loadedContract.RightToUses()
		assert.Len(t, rightToUses, 2)
		assert.True(t, rightToUses[0].ValidFrom().Equal(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)))
		assert.True(t, rightToUses[0].ValidTo().Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
		assert.NotZero(t, rightToUses[0].BillDetailId())
		assert.True(t, rightToUses[1].ValidFrom().Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
		assert.True(t, rightToUses[1].ValidTo().Equal(utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0)))
		assert.Zero(t, rightToUses[1].BillDetailId())
	})

	t.Run("データがない時はnilが返る", func(t *testing.T) {
		r := NewContractRepository()
		// データ取得
		loadedContract, err := r.GetById(-100, db)
		assert.NoError(t, err)
		assert.Nil(t, loadedContract)
	})

	t.Run("使用権がない契約も返る", func(t *testing.T) {
		////// 準備
		r := NewContractRepository()
		// データ登録
		contractEntity := contract.NewContractEntity(
			savedUserId,
			savedProductId,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{},
		)
		savedId, err := r.Create(contractEntity, db)
		assert.NoError(t, err)

		////// 実行
		loadedContract, err := r.GetById(savedId, db)
		assert.NoError(t, err)

		////// 検証
		assert.Equal(t, savedId, loadedContract.Id())
	})

	t.Run("使用権がない(historyには存在する)契約も返る", func(t *testing.T) {
		////// 準備
		r := NewContractRepository()
		// データ登録
		contractEntity := contract.NewContractEntity(
			savedUserId,
			savedProductId,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
				),
			},
		)
		savedId, err := r.Create(contractEntity, db)
		assert.NoError(t, err)
		// 取得して使用権をアーカイブ
		contract, err := r.GetById(savedId, db)
		contract.ArchiveRightToUseByValidTo(utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0))
		err = r.Update(contract, db)
		assert.NoError(t, err)

		////// 実行
		loadedContract, err := r.GetById(savedId, db)
		assert.NoError(t, err)

		////// 検証
		assert.Equal(t, savedId, loadedContract.Id())
	})
}

func TestContractRepository_GetBillingTargetByBillingDate(t *testing.T) {
	db, err := GetConnection()
	assert.NoError(t, err)

	r := NewContractRepository()
	t.Run("請求実行日を渡すとその日以前で請求実行をしていない使用権データを持っているContractがContractId順で全て返る", func(t *testing.T) {
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
		assert.Len(t, actual, 2)
		// 1つめ
		assert.False(t, actual[0].RightToUses()[0].WasBilling())
		assert.Equal(t, rightToUseIds1[0], actual[0].RightToUses()[0].Id())
		assert.True(t, actual[0].RightToUses()[1].WasBilling())
		assert.Equal(t, rightToUseIds1[1], actual[0].RightToUses()[1].Id())
		assert.False(t, actual[0].RightToUses()[2].WasBilling())
		assert.Equal(t, rightToUseIds1[2], actual[0].RightToUses()[2].Id())
		assert.False(t, actual[0].RightToUses()[3].WasBilling())
		assert.Equal(t, rightToUseIds1[3], actual[0].RightToUses()[3].Id())
		// 2つめ
		assert.False(t, actual[1].RightToUses()[0].WasBilling())
		assert.Equal(t, rightToUseIds2[0], actual[1].RightToUses()[0].Id())
		assert.True(t, actual[1].RightToUses()[1].WasBilling())
		assert.Equal(t, rightToUseIds2[1], actual[1].RightToUses()[1].Id())
		assert.False(t, actual[1].RightToUses()[2].WasBilling())
		assert.Equal(t, rightToUseIds2[2], actual[1].RightToUses()[2].Id())
		assert.False(t, actual[1].RightToUses()[3].WasBilling())
		assert.Equal(t, rightToUseIds2[3], actual[1].RightToUses()[3].Id())
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
		assert.Equal(t, rightToUseIds[0], actual[0].RightToUses()[0].Id())
		assert.False(t, actual[0].RightToUses()[0].WasBilling())
		assert.True(t, actual[0].RightToUses()[1].WasBilling())
		assert.False(t, actual[0].RightToUses()[2].WasBilling())
		assert.False(t, actual[0].RightToUses()[3].WasBilling())
	})
}

func deleteTestRightToUseData(executor gorp.SqlExecutor) {
	_, err := executor.Exec(`
DELETE FROM right_to_use_active;
DELETE FROM right_to_use_history;
DELETE FROM discount_apply_contract_updates;
DELETE FROM bill_details;
DELETE FROM right_to_use;
`)
	if err != nil {
		panic("事前データ削除失敗" + err.Error())
	}

	_, err = executor.Exec(`
delete from right_to_use where id in (
    select rtu.id from right_to_use rtu
    left outer join bill_details bd on rtu.id = bd.right_to_use_id
    where bd.id is null
);
`)
	if err != nil {
		panic("事前データ削除失敗")
	}
}

func createBillTestData(db gorp.SqlExecutor) (rightToUseIds []int, billId int) {
	// 商品登録
	productId := createProduct(db)

	// user作成
	// 6/1 ~ 6/30 の使用権（未請求）=> 取得される
	// 6/1 ~ 6/30 の使用権（請求済）
	// 7/1 ~ 7/31 の使用権（未請求）=> 取得される
	// 8/1 ~ 8/31 の使用権（未請求）
	userId := createUser(db)
	savedContract := createContract(
		userId,
		productId,
		utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
		utils.CreateJstTime(2020, 6, 11, 0, 0, 0, 0),
		[]*contract.RightToUseEntity{
			contract.NewRightToUseEntity(
				utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
			),
			contract.NewRightToUseEntity(
				utils.CreateJstTime(2020, 6, 1, 18, 21, 0, 0),
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
			),
			contract.NewRightToUseEntity(
				utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
			),
			contract.NewRightToUseEntity(
				utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 9, 1, 0, 0, 0, 0),
			),
		},
		db,
	)

	// rightToUse1bに対して請求情報を登録しておく（請求済にしておく）
	rightToUses := savedContract.RightToUses()
	billDetailEntity := entities.NewBillingDetailEntity(rightToUses[1].Id(), decimal.NewFromInt(1000))
	billAgg := entities.NewBillingAggregation(utils.CreateJstTime(2020, 7, 1, 12, 0, 0, 0), userId)
	err := billAgg.AddBillDetail(billDetailEntity)
	if err != nil {
		panic("billデータ作成失敗")
	}

	billRep := repositories.NewBillRepository()
	billId, err = billRep.Create(billAgg, db)
	if err != nil {
		panic("billデータ保存失敗")
	}

	return []int{rightToUses[0].Id(), rightToUses[1].Id(), rightToUses[2].Id(), rightToUses[3].Id()}, billId
}

func TestContractRepository_GetRecurTarget(t *testing.T) {
	t.Run("2020/6/1を渡すと使用期間終了日が6/1 ~ 6/5でかつ次の期間の使用権がまだない使用権をもつ契約集約が返る", func(t *testing.T) {
		db, err := GetConnection()
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
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 5, 31, 0, 0, 0, 0),
				),
			},
			tran,
		)
		contractB := createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0),
				),
			},
			tran,
		)
		contractC := createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 6, 5, 0, 0, 0, 0),
				),
			},
			tran,
		)
		_ = createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 6, 5, 0, 0, 0, 0),
				),
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 6, 6, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 7, 5, 0, 0, 0, 0),
				),
			},
			tran,
		)
		_ = createPreparedContractData(
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 4, 1, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 5, 1, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 6, 6, 0, 0, 0, 0),
				),
			},
			tran,
		)
		// 使用権を作成（終了日が5/31, 6/1, 6/5, 6/5（ただし次の使用権データがある） 6/6）
		contractRep := NewContractRepository()

		////// 実行
		expects, err := contractRep.GetRecurTargets(utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0), tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, expects, 2)
		assert.Equal(t, contractB.Id(), expects[0].Id())
		assert.Equal(t, contractC.Id(), expects[1].Id())
	})

	t.Run("更新対象がなければ空スライスが返る", func(t *testing.T) {
		////// 準備
		// 契約データ全削除
		db, err := GetConnection()
		defer db.Db.Close()
		tran, err := db.Begin()
		assert.NoError(t, err)
		_, err = tran.Exec("DELETE FROM discount_apply_contract_updates")
		assert.NoError(t, err)
		_, err = tran.Exec("DELETE FROM bill_details")
		assert.NoError(t, err)
		_, err = tran.Exec("DELETE FROM right_to_use_active")
		assert.NoError(t, err)
		_, err = tran.Exec("DELETE FROM right_to_use_history")
		assert.NoError(t, err)
		_, err = tran.Exec("DELETE FROM right_to_use")
		assert.NoError(t, err)
		_, err = tran.Exec("DELETE FROM contracts")
		assert.NoError(t, err)

		////// 実行
		contractRep := NewContractRepository()
		expects, err := contractRep.GetRecurTargets(utils.CreateJstTime(3020, 6, 1, 0, 0, 0, 0), tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, expects, 0)
	})
}

// 使用権データを作成するのに事前に必要なデータを準備する
func createPreparedContractData(contractDate, billingStartDate time.Time, rightToUses []*contract.RightToUseEntity, executor gorp.SqlExecutor) *contract.ContractEntity {
	// userの作成
	savedUserId := createUser(executor)
	// 商品の作成
	savedProductId := createProduct(executor)
	// 契約の作成
	savedContract := createContract(
		savedUserId,
		savedProductId,
		contractDate,
		billingStartDate,
		rightToUses,
		executor,
	)
	return savedContract
}

func createUser(executor gorp.SqlExecutor) int {
	userEntity, err := user.NewUserIndividualEntity("個人太郎")
	if err != nil {
		panic("userEntity作成失敗")
	}
	userRepository := NewUserRepository()
	savedUserId, err := userRepository.SaveUserIndividual(userEntity, executor)
	if err != nil {
		panic(err.Error())
	}
	return savedUserId
}

func createProduct(executor gorp.SqlExecutor) int {
	productEntity, err := product.NewProductEntity("商品", "1000")
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

func createContract(userId, productId int, contractDate, billingStartDate time.Time, rightToUses []*contract.RightToUseEntity, executor gorp.SqlExecutor) *contract.ContractEntity {
	// 契約の作成
	contractEntity := contract.NewContractEntity(
		userId,
		productId,
		contractDate,
		billingStartDate,
		rightToUses,
	)
	contractRepository := NewContractRepository()
	savedContractId, err := contractRepository.Create(contractEntity, executor)
	if err != nil {
		panic("contractEntity保存失敗")
	}
	// 再読込
	savedContract, err := contractRepository.GetById(savedContractId, executor)
	if err != nil {
		panic("contractEntity取得失敗")
	}
	return savedContract
}

func TestContractRepository_Update(t *testing.T) {
	t.Run("契約エンティティに使用権を追加したものを渡すと、そのデータで更新する", func(t *testing.T) {
		db, err := GetConnection()
		assert.NoError(t, err)
		////// 準備
		contractEntity := createTestContract(db)

		////// 実行
		contractEntity.AddNextTermRightToUses(contract.NewRightToUseEntity(
			utils.CreateJstTime(2020, 11, 5, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 12, 5, 0, 0, 0, 0),
		))
		contractRep := NewContractRepository()
		err = contractRep.Update(contractEntity, db)
		assert.NoError(t, err)

		////// 検証
		actual, err := contractRep.GetById(contractEntity.Id(), db)
		assert.NoError(t, err)
		rightToUses := actual.RightToUses()
		assert.Len(t, rightToUses, 2)
		assert.NotZero(t, rightToUses[0].Id())
		assert.NotZero(t, rightToUses[0].CreatedAt())
		assert.NotZero(t, rightToUses[0].UpdatedAt())
		assert.False(t, rightToUses[0].WasBilling())
		assert.True(t, rightToUses[0].ValidFrom().Equal(utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0)))
		assert.True(t, rightToUses[0].ValidTo().Equal(utils.CreateJstTime(2020, 11, 5, 0, 0, 0, 0)))

		assert.NotZero(t, rightToUses[1].Id())
		assert.NotZero(t, rightToUses[1].CreatedAt())
		assert.NotZero(t, rightToUses[1].UpdatedAt())
		assert.False(t, rightToUses[1].WasBilling())
		assert.True(t, rightToUses[1].ValidFrom().Equal(utils.CreateJstTime(2020, 11, 5, 0, 0, 0, 0)))
		assert.True(t, rightToUses[1].ValidTo().Equal(utils.CreateJstTime(2020, 12, 5, 0, 0, 0, 0)))
	})

	t.Run("アーカイブされた使用権は、historyテーブルへ移動させられる", func(t *testing.T) {
		db, err := GetConnection()
		assert.NoError(t, err)
		////// 準備
		contract := createTestContract(db)
		err = contract.ArchiveRightToUseById(contract.RightToUses()[0].Id())
		assert.NoError(t, err)

		////// 実行
		rep := NewContractRepository()
		err = rep.Update(contract, db)
		assert.NoError(t, err)

		////// 検証（activeテーブルに0件、historyテーブルに1件データがあるはず。）
		sql := `
SELECT
       COUNT(rtu.id)
           FROM
                right_to_use_active rtua
                    INNER JOIN right_to_use rtu ON rtua.right_to_use_id = rtu.id
WHERE rtu.contract_id = $1
;`
		activeCount, err := db.SelectInt(sql, contract.Id())
		assert.NoError(t, err)
		assert.Equal(t, 0, int(activeCount))

		sql = `
SELECT
       COUNT(rtu.id)
           FROM
                right_to_use_history rtuh
                    INNER JOIN right_to_use rtu ON rtuh.right_to_use_id = rtu.id
WHERE rtu.contract_id = $1
;`
		historyCount, err := db.SelectInt(sql, contract.Id())
		assert.NoError(t, err)
		assert.Equal(t, 1, int(historyCount))

	})
}

func createTestContract(db gorp.SqlExecutor) *contract.ContractEntity {
	userID := createUser(db)
	productId := createProduct(db)
	contract := createContract(
		userID,
		productId,
		utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
		[]*contract.RightToUseEntity{
			contract.NewRightToUseEntity(
				utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 11, 5, 0, 0, 0, 0),
			),
		},
		db,
	)
	return contract
}

func TestContractRepository_GetHavingExpiredRightToUseContract(t *testing.T) {
	t.Run("渡した基準日時点で期限が切れているactiveな使用権を持っている契約エンティティを返す", func(t *testing.T) {
		////// 準備（期限が切れている使用権をもってる契約と、持ってない契約を用意する）
		db, err := GetConnection()
		assert.NoError(t, err)
		tran, err := db.Begin()
		assert.NoError(t, err)
		// 事前データを全て削除する
		deleteSql := `
DELETE FROM discount_apply_contract_updates;
DELETE FROM bill_details;
DELETE FROM right_to_use_active;
DELETE FROM right_to_use_history;
DELETE FROM right_to_use;
DELETE FROM contracts;
`
		_, err = db.Exec(deleteSql)
		assert.NoError(t, err)
		// 今回のテスト用データの作成
		userID := createUser(tran)
		productId := createProduct(tran)
		// 期限切れ使用権を持ってる契約
		expiredContract := createContract(
			userID,
			productId,
			utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 9, 5, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
				),
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 11, 5, 0, 0, 0, 0),
				),
			},
			tran,
		)
		// 期限切れ使用権を持ってない契約
		_ = createContract(
			userID,
			productId,
			utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 10, 5, 0, 0, 0, 0),
			[]*contract.RightToUseEntity{
				contract.NewRightToUseEntity(
					utils.CreateJstTime(2020, 11, 5, 0, 0, 0, 0),
					utils.CreateJstTime(2020, 12, 5, 0, 0, 0, 0),
				),
			},
			tran,
		)

		////// 実行
		rep := NewContractRepository()
		actuals, err := rep.GetHavingExpiredRightToUseContractIds(utils.CreateJstTime(2020, 11, 5, 0, 0, 0, 0), tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, actuals, 1)
		assert.Equal(t, expiredContract.Id(), actuals[0])
	})
}
