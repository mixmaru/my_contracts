package application_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v2"
	"testing"
	"time"
)

func TestProductApplicationService_Register(t *testing.T) {
	productEntity := entities.NewProductEntity("商品名", decimal.NewFromFloat(1000))
	returnProductEntity := entities.NewProductEntityWithData(
		100,
		"商品名",
		decimal.NewFromFloat(1000),
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	productRepositoryMock := mock_interfaces.NewMockIProductRepository(ctrl)
	productRepositoryMock.EXPECT().
		Save(
			productEntity,
			gomock.AssignableToTypeOf(&gorp.Transaction{}),
		).Return(returnProductEntity, nil).
		Times(1)

	productApp := NewProductApplicationServiceWithMock(productRepositoryMock)
	dto, _, err := productApp.Register("商品名", decimal.NewFromFloat(1000))
	assert.NoError(t, err)

	assert.Equal(t, 100, dto.Id)
	assert.Equal(t, "商品名", dto.Name)
	assert.True(t, dto.Price.Equal(decimal.NewFromFloat(1000)))
	assert.True(t, dto.CreatedAt.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, dto.UpdatedAt.Equal(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
}

func TestProductApplicationService_Get(t *testing.T) {
	returnProductEntity := entities.NewProductEntityWithData(
		100,
		"商品名",
		decimal.NewFromFloat(1000),
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	)

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
	assert.True(t, dto.Price.Equal(decimal.NewFromFloat(1000)))
	assert.True(t, dto.CreatedAt.Equal(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)))
	assert.True(t, dto.UpdatedAt.Equal(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)))
}
