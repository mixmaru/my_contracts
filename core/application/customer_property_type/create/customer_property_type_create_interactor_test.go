package create

import (
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomerPropertyTypeCreateInteractor_Register(t *testing.T) {
	interactor := NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())

	t.Run("カスタマープロパティ名と型（string or numeric）を渡すとカスタマープロパティデータが作成される", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		request := NewCustomerPropertyTypeCreateUseCaseRequest("性別"+timestampstr, "string")
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		assert.Len(t, response.ValidationError, 0)
		assert.NotZero(t, response.CustomerPropertyTypeDto.Id)
		assert.Equal(t, "性別"+timestampstr, response.CustomerPropertyTypeDto.Name)
		assert.Equal(t, "string", response.CustomerPropertyTypeDto.Type)

		t.Run("既に登録されているプロパティ名だった場合はバリデーションエラーになる", func(t *testing.T) {
			request := NewCustomerPropertyTypeCreateUseCaseRequest("性別"+timestampstr, "string")
			response, err := interactor.Handle(request)
			assert.NoError(t, err)
			expect := map[string][]string{
				"name": []string{
					"既に存在する名前です",
				},
			}

			assert.Len(t, response.ValidationError, 1)
			assert.Equal(t, expect, response.ValidationError)
			assert.Zero(t, response.CustomerPropertyTypeDto.Id)
		})
	})

	t.Run("型にstring or numeric以外の文字がセットされていた場合はバリデーションエラーになる", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		request := NewCustomerPropertyTypeCreateUseCaseRequest("性別"+timestampstr, "hogehoge")
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		expect := map[string][]string{
			"type": []string{
				"stringでもnumericでもありません",
			},
		}

		assert.Len(t, response.ValidationError, 1)
		assert.Equal(t, expect, response.ValidationError)
		assert.Zero(t, response.CustomerPropertyTypeDto.Id)
	})

	t.Run("同時に同じ名前で登録処理が走っても重複エラーが起きない（ファントムリードでのエラーが起きない）", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		request := NewCustomerPropertyTypeCreateUseCaseRequest("同時実行テスト"+timestampstr, "string")

		type ret struct {
			response *CustomerPropertyTypeCreateUseCaseResponse
			error    error
		}
		retCh := make(chan ret)
		for i := 0; i < 50; i++ {
			go func() {
				response, err := interactor.Handle(request)
				retCh <- ret{
					response: response,
					error:    err,
				}
			}()
		}

		// ゴルーチンで並列でたくさん実行してエラーがでないかテストしてる
		var savedResponse *CustomerPropertyTypeCreateUseCaseResponse
		for i := 0; i < 50; i++ {
			select {
			case retVal := <-retCh:
				if retVal.error != nil {
					assert.Failf(t, "エラーが発生した。", retVal.error.Error())
				}
				if savedResponse == nil {
					// 登録実行されてもOK
					if len(retVal.response.ValidationError) == 0 {
						savedResponse = retVal.response
					} else {
						assert.Failf(t, "なぜかバリデーションに失敗してる。", "%v", retVal.response)
					}
				} else {
					// バリデーションエラーになってればOK
					assert.Len(t, retVal.response.ValidationError, 1)
				}
			}
		}
	})
}

//func TestProductApplicationService_registerValidation(t *testing.T) {
//
//t.Run("バリデーションエラーにならない場合はvalidationErrorsは空スライスが返ってくる", func(t *testing.T) {
//	validationErrors, err := createValidation("商品", "1000.01")
//	assert.NoError(t, err)
//	assert.Equal(t, map[string][]string{}, validationErrors)
//})
//
//t.Run("nameが50文字より多い priceがdecimalに変換不可能の場合_バリデーションエラーメッセージが返ってくる", func(t *testing.T) {
//	validationErrors, err := createValidation("1234567890123456789012345678901234567890１２３４５６７８９０1", "aaa")
//	assert.NoError(t, err)
//	expect := map[string][]string{
//		"name": []string{
//			"50文字より多いです",
//		},
//		"price": []string{
//			"数値ではありません",
//		},
//	}
//	assert.Equal(t, expect, validationErrors)
//})
//
//t.Run("nameが空 priceがマイナスだった場合_バリデーションエラーメッセージが返ってくる", func(t *testing.T) {
//	validationErrors, err := createValidation("", "-1000")
//	assert.NoError(t, err)
//	expect := map[string][]string{
//		"name": []string{
//			"空です",
//		},
//		"price": []string{
//			"マイナス値です",
//		},
//	}
//	assert.Equal(t, expect, validationErrors)
//})
//}
