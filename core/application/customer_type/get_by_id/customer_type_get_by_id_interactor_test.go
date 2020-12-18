package get_by_id

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	create2 "github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	"github.com/mixmaru/my_contracts/core/application/customer_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomerTypeGetByIdInteractor_Handle(t *testing.T) {
	t.Run("customerTypeIdを渡すとデータが取得できる", func(t *testing.T) {
		t.Run("存在するIdを渡すとデータが取得できる", func(t *testing.T) {
			////// 準備
			// 事前にCustomerTypeを登録しておく
			preCreateCustomerType, err := preCreateCustomerType()
			assert.NoError(t, err)

			////// 実行
			interactor := NewCustomerTypeGetByIdInteractor(db.NewCustomerTypeRepository(), db.NewCustomerPropertyTypeRepository())
			request := NewCustomerTypeGetByIdUseCaseRequest(preCreateCustomerType.Id)
			response, err := interactor.Handle(request)
			assert.NoError(t, err)

			////// 検証
			assert.Len(t, response.ValidationError, 0)
			assert.Equal(t, preCreateCustomerType, response.CustomerTypeDto)
		})

		t.Run("存在しないIdを渡すとゼロ値が返ってくる", func(t *testing.T) {
			assert.Failf(t, "テストしてない", "テストしてない")
		})
	})
}

func preCreateCustomerType() (customer_type.CustomerTypeDto, error) {
	timestampStr := utils.CreateTimestampString()
	// カスタマープロパティタイプ登録
	preCreatedCustomerPropertyType, err := preCreateCustomerPropertyType()
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}

	// カスタマータイプ登録
	interactor := create.NewCustomerTypeCreateInteractor(db.NewCustomerTypeRepository(), db.NewCustomerPropertyTypeRepository())
	request := create.NewCustomerTypeCreateUseCaseRequest("事前登録カスタマータイプ"+timestampStr, []int{preCreatedCustomerPropertyType.Id})
	response, err := interactor.Handle(request)
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	if len(response.ValidationError) > 0 {
		return customer_type.CustomerTypeDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationError)
	}

	return response.CustomerTypeDto, nil
}

func preCreateCustomerPropertyType() (customer_property_type.CustomerPropertyTypeDto, error) {
	timestampStr := utils.CreateTimestampString()
	interactor := create2.NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())

	request := create2.NewCustomerPropertyTypeCreateUseCaseRequest("事前登録カスタマープロパティタイプ"+timestampStr, "string")
	response, err := interactor.Handle(request)
	if err != nil {
		return customer_property_type.CustomerPropertyTypeDto{}, err
	}
	if len(response.ValidationError) > 0 {
		return customer_property_type.CustomerPropertyTypeDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationError)
	}

	return response.CustomerPropertyTypeDto, nil
}
