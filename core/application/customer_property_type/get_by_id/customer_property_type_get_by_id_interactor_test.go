package get_by_id

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomerPropertyTypeGetByIdInteractor_Register(t *testing.T) {
	interactor := NewCustomerPropertyTypeGetByIdInteractor(db.NewCustomerPropertyTypeRepository())

	t.Run("idを渡すとカスタマープロパティのデータを取得できる", func(t *testing.T) {
		////// 準備
		// カスタマープロパティタイプを登録する
		id, err := preCreate()
		assert.NoError(t, err)

		////// 実行
		request := NewCustomerPropertyTypeGetByIdUseCaseRequest(id)
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, response.ValidationError, 0)
		assert.Equal(t, id, response.CustomerPropertyTypeDto.Id)
	})

	t.Run("存在しないidを渡すとゼロ値で返ってくる", func(t *testing.T) {
		////// 実行
		request := NewCustomerPropertyTypeGetByIdUseCaseRequest(-100)
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, response.ValidationError, 0)
		assert.Zero(t, response.CustomerPropertyTypeDto)
	})
}

func preCreate() (int, error) {
	createInteractor := create.NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())

	timestampstr := utils.CreateTimestampString()
	request := create.NewCustomerPropertyTypeCreateUseCaseRequest("取得テスト"+timestampstr, "string")
	response, err := createInteractor.Handle(request)
	if err != nil {
		return 0, err
	}
	if len(response.ValidationErrors) > 0 {
		return 0, errors.Errorf("バリデーションエラー。%+v", response.ValidationErrors)
	}
	return response.CustomerPropertyTypeDto.Id, nil
}
