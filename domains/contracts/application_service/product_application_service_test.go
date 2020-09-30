package application_service

import (
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestProductApplicationService_Register(t *testing.T) {
	productApp := NewProductApplicationService()

	t.Run("商品名と金額を渡すと商品データが作成される", func(t *testing.T) {
		// 重複しない商品名でテストを行う
		unixNano := time.Now().UnixNano()
		suffix := strconv.FormatInt(unixNano, 10)
		name := "商品" + suffix

		dto, _, err := productApp.Register(name, "1000")
		assert.NoError(t, err)

		assert.NotZero(t, dto.Id)
		assert.Equal(t, name, dto.Name)
		assert.Equal(t, "1000", dto.Price)
		assert.NotZero(t, dto.CreatedAt)
		assert.NotZero(t, dto.UpdatedAt)
	})
}

func TestProductApplicationService_Get(t *testing.T) {
	// 事前にデータ登録（取得テスト用）
	dto := createProduct()
	productApp := NewProductApplicationService()
	t.Run("データがある時はidで取得できる", func(t *testing.T) {
		actual, err := productApp.Get(dto.Id)
		assert.NoError(t, err)

		assert.Equal(t, dto.Id, actual.Id)
		assert.Equal(t, dto.Name, actual.Name)
		assert.Equal(t, dto.Price, actual.Price)
		assert.True(t, actual.CreatedAt.Equal(dto.CreatedAt))
		assert.True(t, actual.UpdatedAt.Equal(dto.UpdatedAt))
	})

	t.Run("データがない時はゼロ値が返る", func(t *testing.T) {
		dto, err := productApp.Get(-100)
		assert.NoError(t, err)

		assert.Zero(t, dto)
	})
}

func TestProductApplicationService_registerValidation(t *testing.T) {
	conn, err := db_connection.GetConnection()
	assert.NoError(t, err)

	productAppService := NewProductApplicationService()

	t.Run("バリデーションエラーにならない場合はvalidationErrorsは空スライスが返ってくる", func(t *testing.T) {
		validationErrors, err := productAppService.registerValidation(utils.CreateUniqProductNameForTest(), "1000.01", conn)
		assert.NoError(t, err)
		assert.Equal(t, map[string][]string{}, validationErrors)
	})

	t.Run("nameが50文字より多い priceがdecimalに変換不可能の場合_バリデーションエラーメッセージが返ってくる", func(t *testing.T) {
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

	t.Run("nameが空 priceがマイナスだった場合_バリデーションエラーメッセージが返ってくる", func(t *testing.T) {
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
