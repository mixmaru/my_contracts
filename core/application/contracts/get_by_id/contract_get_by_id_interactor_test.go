package get_by_id

import (
	contract_create "github.com/mixmaru/my_contracts/core/application/contracts/create"
	"github.com/mixmaru/my_contracts/core/application/products"
	"github.com/mixmaru/my_contracts/core/application/products/create"
	"github.com/mixmaru/my_contracts/core/application/users"
	user_create "github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContractGetByIdInteractor_Handle(t *testing.T) {
	userDto := createUser()
	productDto := createProduct()
	createInteractor := contract_create.NewContractCreateInteractor(db.NewUserRepository(), db.NewProductRepository(), db.NewContractRepository())
	createResponse, err := createInteractor.Handle(contract_create.NewContractCreateUseCaseRequest(userDto.Id, productDto.Id, utils.CreateJstTime(2020, 1, 2, 2, 0, 0, 0)))
	if err != nil || len(createResponse.ValidationErrors) > 0 {
		panic("データ作成失敗")
	}

	t.Run("Idを渡すと対応するデータが取得できる", func(t *testing.T) {
		getByIdinteractor := NewContractGetByIdInteractor(db.NewContractRepository(), db.NewProductRepository(), db.NewUserRepository())
		t.Run("データがある時はデータが取得できる", func(t *testing.T) {
			response, err := getByIdinteractor.Handle(NewContractGetByIdUseCaseRequest(createResponse.ContractDto.Id))
			assert.NoError(t, err)

			contract := response.ContractDto
			assert.Equal(t, createResponse.ContractDto.Id, contract.Id)
			assert.Equal(t, productDto.Id, contract.ProductId)
			assert.Equal(t, userDto.Id, contract.UserId)
			assert.True(t, contract.ContractDate.Equal(utils.CreateJstTime(2020, 1, 2, 2, 0, 0, 0)))
			assert.True(t, contract.BillingStartDate.Equal(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0)))
			assert.True(t, contract.CreatedAt.Equal(createResponse.ContractDto.CreatedAt))
			assert.True(t, contract.UpdatedAt.Equal(createResponse.ContractDto.UpdatedAt))

			product := response.ProductDto
			assert.Equal(t, productDto.Id, product.Id)
			assert.Equal(t, productDto.Name, product.Name)
			assert.Equal(t, "2000", product.Price)
			assert.True(t, product.CreatedAt.Equal(productDto.CreatedAt))
			assert.True(t, product.UpdatedAt.Equal(productDto.UpdatedAt))

			user := response.UserDto
			gotUserDto, ok := user.(users.UserIndividualDto)
			assert.True(t, ok)
			assert.Equal(t, userDto.Id, gotUserDto.Id)
			assert.Equal(t, "個人たろう", gotUserDto.Name)
			assert.True(t, gotUserDto.CreatedAt.Equal(userDto.CreatedAt))
			assert.True(t, gotUserDto.UpdatedAt.Equal(userDto.UpdatedAt))
		})

		t.Run("データがない時はゼロ値が返ってくる", func(t *testing.T) {
			// 実行
			response, err := getByIdinteractor.Handle(NewContractGetByIdUseCaseRequest(-100))
			assert.NoError(t, err)

			// 検証
			assert.Zero(t, response.ContractDto)
			assert.Zero(t, response.ProductDto)
			assert.Nil(t, response.UserDto)
		})
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
	interactor := user_create.NewUserIndividualCreateInteractor(db.NewUserRepository())
	response, err := interactor.Handle(user_create.NewUserIndividualCreateUseCaseRequest("個人たろう"))
	if err != nil || len(response.ValidationErrors) > 0 {
		panic("データ作成失敗")
	}
	return response.UserDto
}
