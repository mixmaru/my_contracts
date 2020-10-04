package domain_service

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/data_mappers"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

func TestContractDomainService_CreateContract(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	userId, productId := createPrepareData(db, t)
	domainService := createContractDomainService()
	t.Run("ユーザーIDと商品IDと契約作成日と課金開始日を渡すと_契約と使用権データを作成してDBに保存し_保存したデータを返す", func(t *testing.T) {
		tran, err := db.Begin()
		assert.NoError(t, err)

		// 実行
		actualContractDto, validErrors, err := domainService.CreateContract(
			userId,
			productId,
			utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0),
			tran,
		)
		assert.NoError(t, err)
		assert.Len(t, validErrors, 0)

		err = tran.Commit()
		assert.NoError(t, err)

		// 検証
		assert.NotZero(t, actualContractDto.Id)
		assert.Equal(t, userId, actualContractDto.UserId)
		assert.Equal(t, productId, actualContractDto.ProductId)
		assert.True(t, actualContractDto.ContractDate.Equal(utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0)))
		assert.True(t, actualContractDto.BillingStartDate.Equal(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0)))
		assert.NotZero(t, actualContractDto.ContractDate)
		assert.NotZero(t, actualContractDto.UpdatedAt)

		// データ保存されているか見ておく
		contractRepository := repositories.NewContractRepository()
		contractEntity, _, _, err := contractRepository.GetById(actualContractDto.Id, db)
		assert.NoError(t, err)
		assert.NotZero(t, contractEntity.Id())

		mapper := data_mappers.RightToUseMapper{}
		err = db.SelectOne(&mapper, "SELECT * FROM right_to_use where contract_id = $1", contractEntity.Id())
		assert.NoError(t, err)
		assert.NotZero(t, mapper.Id)
		assert.Equal(t, contractEntity.Id(), mapper.ContractId)
		assert.True(t, mapper.ValidFrom.Equal(utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0)))
		assert.True(t, mapper.ValidTo.Equal(utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0)))
		assert.NotZero(t, mapper.CreatedAt)
		assert.NotZero(t, mapper.UpdatedAt)
	})
}

func TestNewContractDomainService_calculateBillingStartDate(t *testing.T) {
	app := NewContractDomainService(nil, nil, nil)
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
		rightToUse := entities.NewRightToUseEntityWithData(1,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 2, 1, 0, 0, 0, 0),
			10,
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
			utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0),
		)
		product, err := entities.NewProductEntityWithData(
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

func createPrepareData(db gorp.SqlExecutor, t *testing.T) (userId, productId int) {

	// 事前準備。userを登録しとく
	userEntity, err := entities.NewUserIndividualEntity("個人たろう")
	assert.NoError(t, err)
	userRepository := repositories.NewUserRepository()
	userId, err = userRepository.SaveUserIndividual(userEntity, db)
	assert.NoError(t, err)

	// 事前準備。productを登録しとく
	// 重複しない商品名でテストを行う
	productEntity, err := entities.NewProductEntity("商品", "200")
	assert.NoError(t, err)
	productRepository := repositories.NewProductRepository()
	productId, err = productRepository.Save(productEntity, db)
	assert.NoError(t, err)

	return userId, productId
}

func createContractDomainService() *ContractDomainService {
	return NewContractDomainService(
		repositories.NewContractRepository(),
		repositories.NewUserRepository(),
		repositories.NewProductRepository(),
	)
}
