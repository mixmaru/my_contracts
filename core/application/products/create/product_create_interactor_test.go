package create

import (
	"testing"

	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/stretchr/testify/assert"
)

func TestProductCreateInteractor_Register(t *testing.T) {
	interactor := NewProductCreateInteractor(db.NewProductRepository())

	t.Run("商品名と金額を渡すと商品データが作成される", func(t *testing.T) {
		request := NewProductCreateUseCaseRequest("商品", "1000")
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		assert.Len(t, response.ValidationError, 0)
		assert.NotZero(t, response.ProductDto.Id)
		assert.Equal(t, "商品", response.ProductDto.Name)
		assert.Equal(t, "1000", response.ProductDto.Price)
		assert.NotZero(t, response.ProductDto.CreatedAt)
		assert.NotZero(t, response.ProductDto.UpdatedAt)
	})
}

func TestProductApplicationService_registerValidation(t *testing.T) {

	t.Run("バリデーションエラーにならない場合はvalidationErrorsは空スライスが返ってくる", func(t *testing.T) {
		validationErrors, err := createValidation("商品", "1000.01")
		assert.NoError(t, err)
		assert.Equal(t, map[string][]string{}, validationErrors)
	})

	t.Run("nameが50文字より多い priceがdecimalに変換不可能の場合_バリデーションエラーメッセージが返ってくる", func(t *testing.T) {
		validationErrors, err := createValidation("1234567890123456789012345678901234567890１２３４５６７８９０1", "aaa")
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

	t.Run("nameが空 priceがマイナスだった場合_バリデーションエラーメッセージが返ってくる", func(t *testing.T) {
		validationErrors, err := createValidation("", "-1000")
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
