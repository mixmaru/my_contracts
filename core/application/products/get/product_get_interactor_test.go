package get

import (
	"github.com/mixmaru/my_contracts/core/application/products/create"
	"github.com/mixmaru/my_contracts/core/application/products/dto"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductGetInteractor_Handle(t *testing.T) {
	// 事前にデータ登録（取得テスト用）
	dto := createProduct()
	interactor := NewProductGetInteractor(db.NewProductRepository())
	t.Run("データがある時はidで取得できる", func(t *testing.T) {
		actual, err := interactor.Handle(NewProductGetUseCaseRequest(dto.Id))
		assert.NoError(t, err)

		assert.Equal(t, dto.Id, actual.ProductDto.Id)
		assert.Equal(t, dto.Name, actual.ProductDto.Name)
		assert.Equal(t, dto.Price, actual.ProductDto.Price)
		assert.True(t, actual.ProductDto.CreatedAt.Equal(dto.CreatedAt))
		assert.True(t, actual.ProductDto.UpdatedAt.Equal(dto.UpdatedAt))
	})

	t.Run("データがない時はゼロ値が返る", func(t *testing.T) {
		response, err := interactor.Handle(NewProductGetUseCaseRequest(-100))
		assert.NoError(t, err)

		assert.Zero(t, response.ProductDto)
	})
}

func createProduct() dto.ProductDto {
	createInteractor := create.NewProductCreateInteractor(db.NewProductRepository())
	response, err := createInteractor.Handle(create.NewProductCreateUseCaseRequest("商品", "2000"))
	if err != nil || len(response.ValidationError) > 0 {
		panic("データ作成失敗")
	}

	return response.ProductDto
}
