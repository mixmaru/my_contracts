package create

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProductCreateInteractor_Register(t *testing.T) {
	interactor := NewProductCreateInteractor()

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

//func TestProductCreateInteractor_Get(t *testing.T) {
//	// 事前にデータ登録（取得テスト用）
//	dto := createProduct()
//	productApp := NewProductCreateInteractor()
//	t.Run("データがある時はidで取得できる", func(t *testing.T) {
//		actual, err := productApp.Get(dto.Id)
//		assert.NoError(t, err)
//
//		assert.Equal(t, dto.Id, actual.Id)
//		assert.Equal(t, dto.Name, actual.Name)
//		assert.Equal(t, dto.Price, actual.Price)
//		assert.True(t, actual.CreatedAt.Equal(dto.CreatedAt))
//		assert.True(t, actual.UpdatedAt.Equal(dto.UpdatedAt))
//	})
//
//	t.Run("データがない時はゼロ値が返る", func(t *testing.T) {
//		dto, err := productApp.Get(-100)
//		assert.NoError(t, err)
//
//		assert.Zero(t, dto)
//	})
//}
//
//func TestProductApplicationService_registerValidation(t *testing.T) {
//	conn, err := db_connection.GetConnection()
//	assert.NoError(t, err)
//
//	productAppService := NewProductCreateInteractor()
//
//	t.Run("バリデーションエラーにならない場合はvalidationErrorsは空スライスが返ってくる", func(t *testing.T) {
//		validationErrors, err := productAppService.registerValidation("商品", "1000.01", conn)
//		assert.NoError(t, err)
//		assert.Equal(t, map[string][]string{}, validationErrors)
//	})
//
//	t.Run("nameが50文字より多い priceがdecimalに変換不可能の場合_バリデーションエラーメッセージが返ってくる", func(t *testing.T) {
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
//	t.Run("nameが空 priceがマイナスだった場合_バリデーションエラーメッセージが返ってくる", func(t *testing.T) {
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
