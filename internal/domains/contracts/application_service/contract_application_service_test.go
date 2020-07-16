package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

func TestContractApplicationService_Register(t *testing.T) {
	returnContractEntity, err := entities.NewContractEntityWithData(
		100,
		2,
		3,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	returnProductEntity, err := entities.NewProductEntityWithData(
		3,
		"商品A",
		"2000",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	returnUserEntity, err := entities.NewUserCorporationEntityWithData(
		2,
		"担当太郎",
		"社長次郎",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	contractRepositoryMock := mock_interfaces.NewMockIContractRepository(ctrl)
	contractRepositoryMock.EXPECT().
		Create(
			gomock.AssignableToTypeOf(&entities.ContractEntity{}),
			gomock.AssignableToTypeOf(&gorp.Transaction{}),
		).Return(100, nil).
		Times(1)
	contractRepositoryMock.EXPECT().
		GetById(
			100,
			gomock.AssignableToTypeOf(&gorp.Transaction{}),
		).Return(returnContractEntity, returnProductEntity, returnUserEntity, nil).
		Times(1)

	app := NewContractApplicationServiceWithMock(contractRepositoryMock)
	dto, _, err := app.Register(2, 3)
	assert.NoError(t, err)

	assert.Equal(t, 100, dto.Id)
	assert.Equal(t, 2, dto.UserId)
	assert.Equal(t, 3, dto.ProductId)
	assert.True(t, dto.CreatedAt.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, dto.UpdatedAt.Equal(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func TestContractApplicationService_GetById(t *testing.T) {
	t.Run("データがある時", func(t *testing.T) {
		returnContractEntity, err := entities.NewContractEntityWithData(
			100,
			2,
			3,
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)

		returnProductEntity, err := entities.NewProductEntityWithData(
			3,
			"商品A",
			"2000",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)

		returnUserEntity, err := entities.NewUserCorporationEntityWithData(
			2,
			"担当太郎",
			"社長次郎",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC),
		)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repositoryMock := mock_interfaces.NewMockIContractRepository(ctrl)
		repositoryMock.EXPECT().
			GetById(
				100,
				gomock.Any(),
			).Return(returnContractEntity, returnProductEntity, returnUserEntity, nil).
			Times(1)

		contractApp := NewContractApplicationServiceWithMock(repositoryMock)
		contract, product, user, err := contractApp.GetById(100)
		assert.NoError(t, err)

		assert.Equal(t, 100, contract.Id)
		assert.Equal(t, 3, contract.ProductId)
		assert.Equal(t, 2, contract.UserId)
		assert.True(t, contract.CreatedAt.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, contract.UpdatedAt.Equal(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))

		assert.Equal(t, 3, product.Id)
		assert.Equal(t, "商品A", product.Name)
		assert.Equal(t, "2000", product.Price)
		assert.True(t, product.CreatedAt.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, product.UpdatedAt.Equal(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)))

		userDto, ok := user.(data_transfer_objects.UserCorporationDto)
		assert.True(t, ok)
		assert.Equal(t, 2, userDto.Id)
		assert.Equal(t, "担当太郎", userDto.ContactPersonName)
		assert.Equal(t, "社長次郎", userDto.PresidentName)
		assert.True(t, userDto.CreatedAt.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, userDto.UpdatedAt.Equal(time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)))
	})

	//t.Run("データがない時", func(t *testing.T) {
	//	ctrl := gomock.NewController(t)
	//	defer ctrl.Finish()
	//	productRepositoryMock := mock_interfaces.NewMockIProductRepository(ctrl)
	//	productRepositoryMock.EXPECT().
	//		GetById(
	//			100,
	//			gomock.Any(),
	//		).Return(nil, nil).
	//		Times(1)
	//
	//	productApp := NewProductApplicationServiceWithMock(productRepositoryMock)
	//	dto, err := productApp.Get(100)
	//	assert.NoError(t, err)
	//
	//	assert.Zero(t, dto)
	//})
}

//func TestContractApplicationService_registerValidation(t *testing.T) {
//	// productデータをすべて削除
//	conn, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//	defer conn.Db.Close()
//	_, err = conn.Exec("truncate products cascade")
//	assert.NoError(t, err)
//
//	// 既存データの作成
//	app := NewProductApplicationService()
//	_, validationErrors, err := app.Register("既存商品", "1000")
//	assert.NoError(t, err)
//	assert.Zero(t, validationErrors)
//
//	productAppService := NewProductApplicationService()
//
//	t.Run("エラーなし", func(t *testing.T) {
//		validationErrors, err := productAppService.registerValidation("A商品", "1000.01", conn)
//		assert.NoError(t, err)
//		assert.Equal(t, map[string][]string{}, validationErrors)
//	})
//
//	t.Run("nameが50文字より多い priceがdecimalに変換不可能", func(t *testing.T) {
//		validationErrors, err := productAppService.registerValidation("1234567890123456789012345678901234567890１２３４５６７８９０1", "aaa", conn)
//		assert.NoError(t, err)
//		expect := map[string][]string{
//			"name": []string{
//				"50文字より多いです",
//			},
//			"price": []string{
//				"数値ではありません",
//			},
//		}
//		assert.Equal(t, expect, validationErrors)
//	})
//
//	t.Run("nameがすでに存在する商品名", func(t *testing.T) {
//		validationErrors, err := productAppService.registerValidation("既存商品", "1000", conn)
//		assert.NoError(t, err)
//		expect := map[string][]string{
//			"name": []string{
//				"すでに存在します",
//			},
//		}
//		assert.Equal(t, expect, validationErrors)
//	})
//
//	t.Run("nameが空 priceがマイナス", func(t *testing.T) {
//		validationErrors, err := productAppService.registerValidation("", "-1000", conn)
//		assert.NoError(t, err)
//		expect := map[string][]string{
//			"name": []string{
//				"空です",
//			},
//			"price": []string{
//				"マイナス値です",
//			},
//		}
//		assert.Equal(t, expect, validationErrors)
//	})
//}
