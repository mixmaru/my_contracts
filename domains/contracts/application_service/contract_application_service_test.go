package application_service

import (
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
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

	t.Run("顧客Idと商品IDを契約日時を渡すと課金開始日が当日で契約が作成される", func(t *testing.T) {
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
		assert.True(t, utils.CreateJstTime(2020, 2, 28, 0, 0, 0, 0).Equal(dto.BillingStartDate))
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
			assert.True(t, contract.BillingStartDate.Equal(utils.CreateJstTime(2020, 1, 2, 0, 0, 0, 0)))
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
	productApp := NewProductApplicationService()
	productDto, validErrors, err := productApp.Register("商品", "2000")
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

func TestContractApplicationService_CreateNextRightToUse(t *testing.T) {
	t.Run("渡した実行日から5日以内に期間終了である使用権に対して、次の期間の使用権データを作成して永続化して返却する", func(t *testing.T) {
		// 事前に影響のあるデータを削除しておく（ちょっと広めに削除）
		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		defer db.Db.Close()
		_, err = db.Exec("DELETE FROM right_to_use_active WHERE right_to_use_id IN (SELECT id FROM right_to_use WHERE '2020-05-25' <= valid_to AND valid_to <= '2020-06-02')")
		assert.NoError(t, err)
		_, err = db.Exec("DELETE FROM right_to_use WHERE '2020-05-25' <= valid_to AND valid_to <= '2020-06-02'")
		assert.NoError(t, err)

		////// 準備（2020-05-31が終了日である使用権と2020-05-29が終了日である使用権を作成する）
		user := createUser()
		product := createProduct()
		contractApp := NewContractApplicationService()
		_, validErrors, err := contractApp.Register(user.Id, product.Id, utils.CreateJstTime(2020, 5, 1, 3, 0, 0, 0))
		if err != nil || len(validErrors) > 0 {
			panic("データ作成失敗")
		}
		_, validErrors, err = contractApp.Register(user.Id, product.Id, utils.CreateJstTime(2020, 4, 30, 0, 0, 0, 0))
		if err != nil || len(validErrors) > 0 {
			panic("データ作成失敗")
		}

		////// 実行
		app := NewContractApplicationService()
		actualContracts, err := app.CreateNextRightToUse(utils.CreateJstTime(2020, 5, 28, 0, 10, 0, 0))
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, actualContracts, 2)
		// 1つめ
		recurRightToUse1 := actualContracts[0].RightToUseDtos[1]
		assert.NotZero(t, recurRightToUse1.Id)
		assert.NotZero(t, recurRightToUse1.CreatedAt)
		assert.NotZero(t, recurRightToUse1.UpdatedAt)
		assert.True(t, recurRightToUse1.ValidFrom.Equal(utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0)))
		assert.True(t, recurRightToUse1.ValidTo.Equal(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
		// 2つめ
		recurRightToUse2 := actualContracts[1].RightToUseDtos[1]
		assert.NotZero(t, recurRightToUse2.Id)
		assert.NotZero(t, recurRightToUse2.CreatedAt)
		assert.NotZero(t, recurRightToUse2.UpdatedAt)
		assert.True(t, recurRightToUse2.ValidFrom.Equal(utils.CreateJstTime(2020, 5, 30, 0, 0, 0, 0)))
		assert.True(t, recurRightToUse2.ValidTo.Equal(utils.CreateJstTime(2020, 6, 30, 0, 0, 0, 0)))
	})
}

func TestContractApplicationService_ArchiveExpiredRightToUse(t *testing.T) {
	t.Run("渡した基準日に期限が切れている使用権をアーカイブ処理し、処理した使用権dtoを返す", func(t *testing.T) {
		////// 準備
		// 事前に存在するデータを削除しておく
		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
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

		user := createUser()
		product := createProduct()
		contractApp := NewContractApplicationService()
		contractDto1, validErrors, err := contractApp.Register(
			user.Id,
			product.Id,
			utils.CreateJstTime(2020, 5, 1, 3, 0, 0, 0))
		if err != nil || len(validErrors) > 0 {
			panic("データ作成失敗")
		}

		contractDto2, validErrors, err := contractApp.Register(
			user.Id,
			product.Id,
			utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0))
		if err != nil || len(validErrors) > 0 {
			panic("データ作成失敗")
		}

		_, validErrors, err = contractApp.Register(
			user.Id,
			product.Id,
			utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0))
		if err != nil || len(validErrors) > 0 {
			panic("データ作成失敗")
		}

		////// 実行
		app := NewContractApplicationService()
		dtos, err := app.ArchiveExpiredRightToUse(utils.CreateJstTime(2020, 6, 2, 0, 0, 0, 0))

		////// 検証
		assert.Len(t, dtos, 2)
		assert.Equal(t, contractDto1.RightToUseDtos[0], dtos[0])
		assert.Equal(t, contractDto2.RightToUseDtos[1], dtos[1])
	})
}
