package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestContractApplicationService_Register(t *testing.T) {
	conn, err := db_connection.GetConnection()
	assert.NoError(t, err)
	tran, err := conn.Begin()
	assert.NoError(t, err)

	savedProductDto := createProduct()
	savedUserDto := createUser()

	err = tran.Commit()
	assert.NoError(t, err)

	app := NewContractApplicationService()

	t.Run("顧客Idと商品IDを契約日時を渡すと課金開始日が翌日で契約が作成される", func(t *testing.T) {
		// 実行
		contractDateTime := utils.CreateJstTime(2020, 2, 28, 23, 0, 0, 0)
		dto, validErrors, err := app.Register(savedUserDto.Id, savedProductDto.Id, contractDateTime)
		assert.NoError(t, err)
		assert.Len(t, validErrors, 0)

		// 検証
		assert.NotZero(t, dto.Id)
		assert.Equal(t, savedUserDto.Id, dto.UserId)
		assert.Equal(t, savedProductDto.Id, dto.ProductId)
		assert.True(t, contractDateTime.Equal(dto.ContractDate))
		assert.True(t, utils.CreateJstTime(2020, 2, 29, 0, 0, 0, 0).Equal(dto.BillingStartDate))
		assert.NotZero(t, dto.CreatedAt)
		assert.NotZero(t, dto.UpdatedAt)
	})

	t.Run("指定されたUserが存在しない時_validationErrorsにエラーメッセージが返ってくる", func(t *testing.T) {
		dto, validationErrors, err := app.Register(-100, savedProductDto.Id, time.Now())
		assert.NoError(t, err)
		assert.Len(t, validationErrors, 1)
		assert.Len(t, validationErrors["user_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["user_id"][0])
		assert.Zero(t, dto)
	})

	t.Run("指定されたProductが存在しない時_validationErrorsにエラーメッセージが返ってくる", func(t *testing.T) {
		dto, validationErrors, err := app.Register(savedUserDto.Id, -100, time.Now())
		assert.NoError(t, err)
		assert.Len(t, validationErrors, 1)
		assert.Len(t, validationErrors["product_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["product_id"][0])
		assert.Zero(t, dto)
	})

	t.Run("指定されたProductもuserも存在しない時_validationErrorsに両方を示すエラーメッセージが返ってくる", func(t *testing.T) {
		dto, validationErrors, err := app.Register(-1000, -100, time.Now())
		assert.NoError(t, err)
		assert.Len(t, validationErrors, 2)
		assert.Len(t, validationErrors["user_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["user_id"][0])
		assert.Len(t, validationErrors["product_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["product_id"][0])
		assert.Zero(t, dto)
	})
}

func TestContractApplicationService_GetById(t *testing.T) {
	userDto := createUser()
	productDto := createProduct()
	contractApp := NewContractApplicationService()
	contractDto, validErrors, err := contractApp.Register(userDto.Id, productDto.Id, utils.CreateJstTime(2020, 1, 2, 2, 0, 0, 0))
	if err != nil || len(validErrors) > 0 {
		panic("データ作成失敗")
	}

	t.Run("Idを渡すと対応するデータが取得できる", func(t *testing.T) {
		t.Run("データがある時はデータが取得できる", func(t *testing.T) {
			contract, product, user, err := contractApp.GetById(contractDto.Id)
			assert.NoError(t, err)

			assert.Equal(t, contractDto.Id, contract.Id)
			assert.Equal(t, productDto.Id, contract.ProductId)
			assert.Equal(t, userDto.Id, contract.UserId)
			assert.True(t, contract.ContractDate.Equal(utils.CreateJstTime(2020, 1, 2, 2, 0, 0, 0)))
			assert.True(t, contract.BillingStartDate.Equal(utils.CreateJstTime(2020, 1, 3, 0, 0, 0, 0)))
			assert.True(t, contract.CreatedAt.Equal(contractDto.CreatedAt))
			assert.True(t, contract.UpdatedAt.Equal(contractDto.UpdatedAt))

			assert.Equal(t, productDto.Id, product.Id)
			assert.Equal(t, productDto.Name, product.Name)
			assert.Equal(t, "2000", product.Price)
			assert.True(t, product.CreatedAt.Equal(productDto.CreatedAt))
			assert.True(t, product.UpdatedAt.Equal(productDto.UpdatedAt))

			gotUserDto, ok := user.(data_transfer_objects.UserIndividualDto)
			assert.True(t, ok)
			assert.Equal(t, userDto.Id, gotUserDto.Id)
			assert.Equal(t, "個人たろう", gotUserDto.Name)
			assert.True(t, gotUserDto.CreatedAt.Equal(userDto.CreatedAt))
			assert.True(t, gotUserDto.UpdatedAt.Equal(userDto.UpdatedAt))
		})

		t.Run("データがない時はゼロ値が返ってくる", func(t *testing.T) {
			// 実行
			contract, product, user, err := contractApp.GetById(-100)
			assert.NoError(t, err)

			// 検証
			assert.Zero(t, contract)
			assert.Zero(t, product)
			assert.Nil(t, user)
		})
	})
}

func createProduct() data_transfer_objects.ProductDto {
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 10)
	name := "商品" + suffix

	productApp := NewProductApplicationService()
	productDto, validErrors, err := productApp.Register(name, "2000")
	if err != nil || len(validErrors) > 0 {
		panic("データ作成失敗")
	}

	return productDto
}

func createUser() data_transfer_objects.UserIndividualDto {
	userApp := NewUserApplicationService()
	dto, validErrors, err := userApp.RegisterUserIndividual("個人たろう")
	if err != nil || len(validErrors) > 0 {
		panic("データ作成失敗")
	}
	return dto
}
