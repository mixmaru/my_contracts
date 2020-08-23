package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
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
	// 事前に31000円の商品を登録
	productEntity := createProduct("31000")

	// ドメインサービスインスタンス化
	billingDS := NewBillingCalculatorDomainService(repositories.NewProductRepository())

	// db接続作成
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)

	t.Run("契約初月", func(t *testing.T) {
		t.Run("31日ある月", func(t *testing.T) {
			// テスト用契約を新規作成する
			contract, err := entities.NewContractEntityWithData(
				1,
				2,
				productEntity.Id(),
				utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 1, 16, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
			)
			assert.NoError(t, err)

			t.Run("課金開始日の翌月同日-1より前の日を渡すと_日割り料金が返る", func(t *testing.T) {
				billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 2, 10, 15, 0, 0, 0), db)
				assert.NoError(t, err)
				assert.Equal(t, "26000", billingAmount.String())
			})
			t.Run("課金開始日の翌月同日-1日を渡すと_まるまる1月分の料金が返る", func(t *testing.T) {
				billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 2, 15, 15, 0, 0, 0), db)
				assert.NoError(t, err)
				assert.Equal(t, "31000", billingAmount.String())
			})
		})

		t.Run("30日ある月", func(t *testing.T) {
			// テスト用契約を新規作成する
			contract, err := entities.NewContractEntityWithData(
				1,
				2,
				productEntity.Id(),
				utils.CreateJstTime(2020, 4, 15, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 4, 16, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 4, 15, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 4, 15, 0, 0, 0, 0),
			)
			assert.NoError(t, err)
			t.Run("課金開始日の翌月同日-1より前の日を渡すと_日割り料金が返る", func(t *testing.T) {
				billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 5, 10, 15, 0, 0, 0), db)
				assert.NoError(t, err)
				assert.Equal(t, "25833", billingAmount.String())
			})
			t.Run("課金開始日の翌月同日-1日を渡すと_まるまる1月分の料金が返る", func(t *testing.T) {
				billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 5, 15, 15, 0, 0, 0), db)
				assert.NoError(t, err)
				assert.Equal(t, "31000", billingAmount.String())
			})
		})

		t.Run("29日ある月", func(t *testing.T) {
			// テスト用契約を新規作成する
			contract, err := entities.NewContractEntityWithData(
				1,
				2,
				productEntity.Id(),
				utils.CreateJstTime(2020, 2, 15, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 2, 16, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 2, 15, 0, 0, 0, 0),
				utils.CreateJstTime(2020, 2, 15, 0, 0, 0, 0),
			)
			assert.NoError(t, err)
			t.Run("課金開始日の翌月同日-1より前の日を渡すと_日割り料金が返る", func(t *testing.T) {
				billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 3, 10, 15, 0, 0, 0), db)
				assert.NoError(t, err)
				assert.Equal(t, "25655", billingAmount.String())
			})
			t.Run("課金開始日の翌月同日-1日を渡すと_まるまる1月分の料金が返る", func(t *testing.T) {
				billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 3, 15, 15, 0, 0, 0), db)
				assert.NoError(t, err)
				assert.Equal(t, "31000", billingAmount.String())
			})
		})
	})

	t.Run("契約翌月以降", func(t *testing.T) {
		// テスト用契約を新規作成する
		contract, err := entities.NewContractEntityWithData(
			1,
			2,
			productEntity.Id(),
			utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 16, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
		)
		assert.NoError(t, err)

		t.Run("課金開始日の翌月同日-1日より後_翌翌月同日-1の日を渡すと_まるまる1月分の料金が返る", func(t *testing.T) {
			billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 3, 15, 15, 0, 0, 0), db)
			assert.NoError(t, err)
			assert.Equal(t, "31000", billingAmount.String())
		})

		t.Run("課金開始日の翌月同日-1日より後_翌翌月同日-1の日より前の日を渡すと_日割り料金が返る", func(t *testing.T) {
			billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 3, 10, 15, 0, 0, 0), db)
			assert.NoError(t, err)
			assert.Equal(t, "25655", billingAmount.String())
		})
	})
}
