package create

import (
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/products/create"
	"github.com/mixmaru/my_contracts/core/application/users"
	create2 "github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestContractCreateInteractor_Handle(t *testing.T) {
	conn, err := db.GetConnection()
	assert.NoError(t, err)
	tran, err := conn.Begin()
	assert.NoError(t, err)

	savedProductDto := createProduct()
	savedUserDto := createUser()

	err = tran.Commit()
	assert.NoError(t, err)

	interactor := NewContractCreateInteractor(db.NewUserRepository(), db.NewProductRepository(), db.NewContractRepository())

	t.Run("顧客Idと商品IDを契約日時を渡すと課金開始日が当日で契約が作成される", func(t *testing.T) {
		// 実行
		contractDateTime := utils.CreateJstTime(2020, 2, 28, 23, 0, 0, 0)
		response, err := interactor.Handle(NewContractCreateUseCaseRequest(savedUserDto.Id, savedProductDto.Id, contractDateTime))
		assert.NoError(t, err)
		assert.Len(t, response.ValidationErrors, 0)

		// 検証
		dto := response.ContractDto
		assert.NotZero(t, dto.Id)
		assert.Equal(t, savedUserDto.Id, dto.UserId)
		assert.Equal(t, savedProductDto.Id, dto.ProductId)
		assert.True(t, contractDateTime.Equal(dto.ContractDate))
		assert.True(t, utils.CreateJstTime(2020, 2, 28, 0, 0, 0, 0).Equal(dto.BillingStartDate))
		assert.NotZero(t, dto.CreatedAt)
		assert.NotZero(t, dto.UpdatedAt)
	})

	t.Run("指定されたUserが存在しない時_validationErrorsにエラーメッセージが返ってくる", func(t *testing.T) {
		response, err := interactor.Handle(NewContractCreateUseCaseRequest(-100, savedProductDto.Id, time.Now()))
		assert.NoError(t, err)
		validationErrors := response.ValidationErrors
		assert.Len(t, validationErrors, 1)
		assert.Len(t, validationErrors["user_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["user_id"][0])
		assert.Zero(t, response.ContractDto)
	})

	t.Run("指定されたProductが存在しない時_validationErrorsにエラーメッセージが返ってくる", func(t *testing.T) {
		response, err := interactor.Handle(NewContractCreateUseCaseRequest(savedUserDto.Id, -100, time.Now()))
		assert.NoError(t, err)
		validationErrors := response.ValidationErrors
		assert.Len(t, validationErrors, 1)
		assert.Len(t, validationErrors["product_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["product_id"][0])
		assert.Zero(t, response.ContractDto)
	})

	t.Run("指定されたProductもuserも存在しない時_validationErrorsに両方を示すエラーメッセージが返ってくる", func(t *testing.T) {
		response, err := interactor.Handle(NewContractCreateUseCaseRequest(-1000, -100, time.Now()))
		assert.NoError(t, err)
		validationErrors := response.ValidationErrors
		assert.Len(t, validationErrors, 2)
		assert.Len(t, validationErrors["user_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["user_id"][0])
		assert.Len(t, validationErrors["product_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["product_id"][0])
		assert.Zero(t, response.ContractDto)
	})
}

func createProduct() products.ProductDto {
	interactor := create.NewProductCreateInteractor(db.NewProductRepository())
	response, err := interactor.Handle(create.NewProductCreateUseCaseRequest("商品", "2000"))
	if err != nil || len(response.ValidationError) > 0 {
		panic("データ作成失敗")
	}

	return response.ProductDto
}

func createUser() users.UserIndividualDto {
	interactor := create2.NewUserIndividualCreateInteractor(db.NewUserRepository())
	response, err := interactor.Handle(create2.NewUserIndividualCreateUseCaseRequest("個人たろう"))
	if err != nil || len(response.ValidationErrors) > 0 {
		panic("データ作成失敗")
	}
	return response.UserDto
}
