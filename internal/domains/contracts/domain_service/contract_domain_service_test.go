package domain_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestContractDomainService_CreateContract(t *testing.T) {
	// 事前準備。userと商品を登録しておく
	userApp := application_service.NewUserApplicationService()
	user, validErrors, err := userApp.RegisterUserIndividual("個人太郎")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)

	productApp := application_service.NewProductApplicationService()
	// 重複しない商品名でテストを行う
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 10)
	name := "商品" + suffix
	product, validErrors, err := productApp.Register(name, "200")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)

	db, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer db.Db.Close()

	domainService := NewContractDomainService(repositories.NewContractRepository(), repositories.NewUserRepository(), repositories.NewProductRepository())
	t.Run("ユーザーIDと商品IDと契約作成日と課金開始日を渡すと_契約と使用権データを作成してDBに保存し_保存したデータを返す", func(t *testing.T) {
		tran, err := db.Begin()
		assert.NoError(t, err)

		actualContractDto, validErrors, err := domainService.CreateContract(
			user.Id,
			product.Id,
			utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0),
			tran,
		)
		assert.NoError(t, err)
		assert.Len(t, validErrors, 0)

		assert.NotZero(t, actualContractDto.Id)
		assert.Equal(t, user.Id, actualContractDto.UserId)
		assert.Equal(t, product.Id, actualContractDto.ProductId)
		assert.True(t, actualContractDto.ContractDate.Equal(utils.CreateJstTime(2020, 1, 1, 15, 0, 0, 0)))
		assert.True(t, actualContractDto.BillingStartDate.Equal(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0)))
		assert.NotZero(t, actualContractDto.ContractDate)
		assert.NotZero(t, actualContractDto.UpdatedAt)
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
