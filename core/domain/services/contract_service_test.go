package services

import (
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"github.com/mixmaru/my_contracts/core/domain/models/product"
	"github.com/mixmaru/my_contracts/core/domain/models/user"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

func TestContractDomainService_CreateContract(t *testing.T) {
	db, err := db.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	userId, productId := createPrepareData(db, t)
	domainService := createContractDomainService()
	t.Run("ユーザーIDと商品IDと契約作成日と課金開始日を渡すと_契約集約を作成して返す", func(t *testing.T) {
		// 実行
		actualContract, validErrors, err := domainService.CreateContract(
			userId,
			productId,
			utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0),
			db,
		)
		assert.NoError(t, err)
		assert.Len(t, validErrors, 0)

		// 検証
		assert.Zero(t, actualContract.Id())
		assert.Zero(t, actualContract.CreatedAt())
		assert.Zero(t, actualContract.UpdatedAt())
		assert.Equal(t, userId, actualContract.UserId())
		assert.Equal(t, productId, actualContract.ProductId())
		assert.True(t, actualContract.ContractDate().Equal(utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0)))
		assert.True(t, actualContract.BillingStartDate().Equal(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)))
	})
}

func TestNewContractDomainService_calculateBillingStartDate(t *testing.T) {
	app := NewContractDomainService(nil, nil)
	t.Run("契約日と無料期間とタイムゾーンを渡すと_課金開始日が返ってくる", func(t *testing.T) {
		t.Run("JSTで渡すと_JSTで0時0分で返ってくる", func(t *testing.T) {
			expect := utils.CreateJstTime(2020, 1, 11, 0, 0, 0, 0)
			actual := app.calculateBillingStartDate(utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0), 10, utils.CreateJstLocation())
			assert.True(t, expect.Equal(actual))
		})
		t.Run("契約開始日をJSTで渡し_locale引数をUTCで渡すと_UTCで0時0分で返ってくる", func(t *testing.T) {
			expect := time.Date(2020, 1, 11, 0, 0, 0, 0, time.UTC)
			actual := app.calculateBillingStartDate(utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0), 10, time.UTC)
			assert.True(t, expect.Equal(actual))
		})
	})
}

func TestNewContractDomainService_CreateNextTermRightToUse(t *testing.T) {
	t.Run("使用権と商品を渡すと次の期間の使用権が作成され返却される", func(t *testing.T) {
		////// 準備
		// 使用権と商品を用意
		rightToUse := contract.NewRightToUseEntityWithData(1,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
			10,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		)
		product, err := product.NewProductEntityWithData(
			10,
			"1ヶ月商品",
			"2000",
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		)
		assert.NoError(t, err)

		////// 実行
		actual, err := CreateNextTermRightToUse(rightToUse, product)
		assert.NoError(t, err)

		////// 検証
		assert.Zero(t, actual.Id())
		assert.Zero(t, actual.CreatedAt())
		assert.Zero(t, actual.UpdatedAt())
		assert.False(t, actual.WasBilling())
		assert.True(t, actual.ValidFrom().Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
		assert.True(t, actual.ValidTo().Equal(utils.CreateJstTime(2020, 3, 1, 0, 0, 0, 0)))
	})
}

func createPrepareData(executor gorp.SqlExecutor, t *testing.T) (userId, productId int) {

	// 事前準備。userを登録しとく
	userEntity, err := user.NewUserIndividualEntity("個人たろう")
	assert.NoError(t, err)
	userRepository := db.NewUserRepository()
	userId, err = userRepository.SaveUserIndividual(userEntity, executor)
	assert.NoError(t, err)

	// 事前準備。productを登録しとく
	// 重複しない商品名でテストを行う
	productEntity, err := product.NewProductEntity("商品", "200")
	assert.NoError(t, err)
	productRepository := db.NewProductRepository()
	productId, err = productRepository.Save(productEntity, executor)
	assert.NoError(t, err)

	return userId, productId
}

func createContractDomainService() *ContractDomainService {
	return NewContractDomainService(
		db.NewUserRepository(),
		db.NewProductRepository(),
	)
}
