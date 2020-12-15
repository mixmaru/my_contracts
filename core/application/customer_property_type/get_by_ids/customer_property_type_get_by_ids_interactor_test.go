package create

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCustomerPropertyTypeGetByIdsInteractor_Register(t *testing.T) {
	interactor := NewCustomerPropertyTypeGetByIdsInteractor(db.NewCustomerPropertyTypeRepository())

	t.Run("idリストを渡すとカスタマープロパティの一覧を取得できる(id順)", func(t *testing.T) {
		////// 準備
		// 2つカスタマープロパティタイプを登録する
		ids, err := preCreate()
		assert.NoError(t, err)

		////// 実行
		request := NewCustomerPropertyTypeGetByIdsUseCaseRequest(ids)
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, response.ValidationError, 0)
		assert.Len(t, response.CustomerPropertyTypeDtos, 2)
		assert.Equal(t, ids[0], response.CustomerPropertyTypeDtos[0].Id)
		assert.Equal(t, ids[1], response.CustomerPropertyTypeDtos[1].Id)
	})
}

func preCreate() ([]int, error) {
	retIds := []int{}

	createInteractor := create.NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())

	timestampstr := utils.CreateTimestampString()

	for i := 0; i < 2; i++ {
		request := create.NewCustomerPropertyTypeCreateUseCaseRequest("取得テスト"+strconv.Itoa(i+1)+timestampstr, "string")
		response, err := createInteractor.Handle(request)
		if err != nil {
			return nil, err
		}
		if len(response.ValidationError) > 0 {
			return nil, errors.Errorf("バリデーションエラー。%+v", response.ValidationError)
		}
		retIds = append(retIds, response.CustomerPropertyTypeDto.Id)
	}
	return retIds, nil
}
