package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/data_mappers"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestContractDomainService_CreateContract(t *testing.T) {
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	// 事前準備。userを登録しとく
	userEntity, err := entities.NewUserIndividualEntity("個人たろう")
	assert.NoError(t, err)
	userRepository := repositories.NewUserRepository()
	userId, err := userRepository.SaveUserIndividual(userEntity, db)
	assert.NoError(t, err)

	// 事前準備。productを登録しとく
	//productApp := application_service.NewProductApplicationService()
	// 重複しない商品名でテストを行う
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 10)
	name := "商品" + suffix
	productEntity, err := entities.NewProductEntity(name, "200")
	assert.NoError(t, err)
	productRepository := repositories.NewProductRepository()
	productId, err := productRepository.Save(productEntity, db)
	assert.NoError(t, err)

	contractRepository := repositories.NewContractRepository()
	rightToUseRepository := repositories.NewRightToUseRepository()

	domainService := NewContractDomainService(contractRepository, userRepository, productRepository, rightToUseRepository)
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
		assert.True(t, actualContractDto.BillingStartDate.Equal(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0)))
		assert.NotZero(t, actualContractDto.ContractDate)
		assert.NotZero(t, actualContractDto.UpdatedAt)

		// データ保存されているか見ておく
		contractEntity, _, _, err := contractRepository.GetById(actualContractDto.Id, db)
		assert.NoError(t, err)
		assert.NotZero(t, contractEntity.Id())

		mapper := data_mappers.RightToUseMapper{}
		err = db.SelectOne(&mapper, "SELECT * FROM right_to_use where contract_id = $1", contractEntity.Id())
		assert.NoError(t, err)
		assert.NotZero(t, mapper.Id)
		assert.Equal(t, contractEntity.Id(), mapper.ContractId)
		assert.True(t, mapper.ValidFrom.Equal(utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0)))
		assert.True(t, mapper.ValidTo.Equal(utils.CreateJstTime(2020, 2, 2, 0, 0, 0, 0)))
		assert.NotZero(t, mapper.CreatedAt)
		assert.NotZero(t, mapper.UpdatedAt)
	})
}

func TestNewContractDomainService_calculateBillingStartDate(t *testing.T) {
	app := NewContractDomainService(nil, nil, nil, nil)
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
