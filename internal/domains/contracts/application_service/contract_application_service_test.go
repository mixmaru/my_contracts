package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestContractApplicationService_Register(t *testing.T) {
	conn, err := db_connection.GetConnection()
	assert.NoError(t, err)
	tran, err := conn.Begin()
	assert.NoError(t, err)

	productRep := repositories.NewProductRepository()

	// productがあればそれを使用する。なければ登録。同名商品は登録できないため
	product, err := productRep.GetByName("バリデーションテスト商品", tran)
	assert.NoError(t, err)
	savedProductId := 0
	if product == nil {
		product, err = entities.NewProductEntity("バリデーションテスト商品", "2000")
		assert.NoError(t, err)
		savedProductId, err = productRep.Save(product, tran)
		assert.NoError(t, err)
	} else {
		savedProductId = product.Id()
	}

	// userを新規登録
	userRep := repositories.NewUserRepository()
	user, err := entities.NewUserIndividualEntity("個人たろう")
	assert.NoError(t, err)
	savedUserId, err := userRep.SaveUserIndividual(user, tran)
	assert.NoError(t, err)

	err = tran.Commit()
	assert.NoError(t, err)

	app := NewContractApplicationService()

	t.Run("エラーなし", func(t *testing.T) {
		dto, validErrors, err := app.Register(savedUserId, savedProductId)
		assert.NoError(t, err)
		assert.Len(t, validErrors, 0)

		assert.NotZero(t, dto.Id)
		assert.Equal(t, savedUserId, dto.UserId)
		assert.Equal(t, savedProductId, dto.ProductId)
		assert.NotZero(t, dto.CreatedAt)
		assert.NotZero(t, dto.UpdatedAt)
	})

	t.Run("指定されたUserが存在しない", func(t *testing.T) {
		dto, validationErrors, err := app.Register(-100, savedProductId)
		assert.NoError(t, err)
		assert.Len(t, validationErrors, 1)
		assert.Len(t, validationErrors["user_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["user_id"][0])
		assert.Zero(t, dto)
	})

	t.Run("指定されたProductが存在しない", func(t *testing.T) {
		dto, validationErrors, err := app.Register(savedUserId, -100)
		assert.NoError(t, err)
		assert.Len(t, validationErrors, 1)
		assert.Len(t, validationErrors["product_id"], 1)
		assert.Equal(t, "存在しません", validationErrors["product_id"][0])
		assert.Zero(t, dto)
	})

	t.Run("指定されたProductもuserも存在しない", func(t *testing.T) {
		dto, validationErrors, err := app.Register(-1000, -100)
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

	t.Run("データがない時", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		repositoryMock := mock_interfaces.NewMockIContractRepository(ctrl)
		repositoryMock.EXPECT().
			GetById(
				100,
				gomock.Any(),
			).Return(nil, nil, nil, nil).
			Times(1)

		contractApp := NewContractApplicationServiceWithMock(repositoryMock)
		contract, product, user, err := contractApp.GetById(100)
		assert.NoError(t, err)

		assert.Zero(t, contract)
		assert.Zero(t, product)
		assert.Nil(t, user)
	})
}
