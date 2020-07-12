package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

func TestProductApplicationService_Register(t *testing.T) {
	returnProductEntity, err := entities.NewProductEntityWithData(
		100,
		"商品名",
		"1000",
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	productRepositoryMock := mock_interfaces.NewMockIProductRepository(ctrl)
	productRepositoryMock.EXPECT().
		Save(
			gomock.AssignableToTypeOf(&entities.ProductEntity{}),
			gomock.AssignableToTypeOf(&gorp.Transaction{}),
		).Return(100, nil).
		Times(1)
	productRepositoryMock.EXPECT().
		GetByName(
			"商品名",
			gomock.AssignableToTypeOf(&gorp.Transaction{}),
		).Return(nil, nil).
		Times(1)
	productRepositoryMock.EXPECT().
		GetById(
			100,
			gomock.AssignableToTypeOf(&gorp.Transaction{}),
		).Return(returnProductEntity, nil).
		Times(1)

	productApp := NewProductApplicationServiceWithMock(productRepositoryMock)
	dto, _, err := productApp.Register("商品名", "1000")
	assert.NoError(t, err)

	assert.Equal(t, 100, dto.Id)
	assert.Equal(t, "商品名", dto.Name)
	assert.Equal(t, "1000", dto.Price)
	assert.True(t, dto.CreatedAt.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, dto.UpdatedAt.Equal(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func TestProductApplicationService_Get(t *testing.T) {
	t.Run("データがある時", func(t *testing.T) {
		returnProductEntity, err := entities.NewProductEntityWithData(
			100,
			"商品名",
			"1000",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		productRepositoryMock := mock_interfaces.NewMockIProductRepository(ctrl)
		productRepositoryMock.EXPECT().
			GetById(
				100,
				gomock.Any(),
			).Return(returnProductEntity, nil).
			Times(1)

		productApp := NewProductApplicationServiceWithMock(productRepositoryMock)
		dto, err := productApp.Get(100)
		assert.NoError(t, err)

		assert.Equal(t, 100, dto.Id)
		assert.Equal(t, "商品名", dto.Name)
		assert.Equal(t, "1000", dto.Price)
		assert.True(t, dto.CreatedAt.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
		assert.True(t, dto.UpdatedAt.Equal(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
	})

	t.Run("データがない時", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		productRepositoryMock := mock_interfaces.NewMockIProductRepository(ctrl)
		productRepositoryMock.EXPECT().
			GetById(
				100,
				gomock.Any(),
			).Return(nil, nil).
			Times(1)

		productApp := NewProductApplicationServiceWithMock(productRepositoryMock)
		dto, err := productApp.Get(100)
		assert.NoError(t, err)

		assert.Zero(t, dto)
	})
}

func TestProductApplicationService_registerValidation(t *testing.T) {
	// productデータをすべて削除
	conn, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer conn.Db.Close()
	_, err = conn.Exec("truncate products cascade")
	assert.NoError(t, err)

	// 既存データの作成
	app := NewProductApplicationService()
	_, validationErrors, err := app.Register("既存商品", "1000")
	assert.NoError(t, err)
	assert.Zero(t, validationErrors)

	productAppService := NewProductApplicationService()

	t.Run("エラーなし", func(t *testing.T) {
		validationErrors, err := productAppService.registerValidation("A商品", "1000.01", conn)
		assert.NoError(t, err)
		assert.Equal(t, map[string][]string{}, validationErrors)
	})

	t.Run("nameが50文字より多い priceがdecimalに変換不可能", func(t *testing.T) {
		validationErrors, err := productAppService.registerValidation("1234567890123456789012345678901234567890１２３４５６７８９０1", "aaa", conn)
		assert.NoError(t, err)
		expect := map[string][]string{
			"name": []string{
				"50文字より多いです",
			},
			"price": []string{
				"数値ではありません",
			},
		}
		assert.Equal(t, expect, validationErrors)
	})

	t.Run("nameがすでに存在する商品名", func(t *testing.T) {
		validationErrors, err := productAppService.registerValidation("既存商品", "1000", conn)
		assert.NoError(t, err)
		expect := map[string][]string{
			"name": []string{
				"すでに存在します",
			},
		}
		assert.Equal(t, expect, validationErrors)
	})

	t.Run("nameが空 priceがマイナス", func(t *testing.T) {
		validationErrors, err := productAppService.registerValidation("", "-1000", conn)
		assert.NoError(t, err)
		expect := map[string][]string{
			"name": []string{
				"空です",
			},
			"price": []string{
				"マイナス値です",
			},
		}
		assert.Equal(t, expect, validationErrors)
	})
}
