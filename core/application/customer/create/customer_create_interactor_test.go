package create

import (
	"github.com/mixmaru/my_contracts/core/application/customer"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	create2 "github.com/mixmaru/my_contracts/core/application/customer_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomerCreateInteractor_Handle(t *testing.T) {
	////// 事前準備
	// カスタマータイプの登録
	customerType, err := preCreateCustomerType()
	assert.NoError(t, err)

	t.Run("リクエストを渡すとカスタマー登録ができる", func(t *testing.T) {
		timestampStr := utils.CreateTimestampString()

		////// 実行
		request := NewCustomerCreateUseCaseRequest(
			"山田"+timestampStr,
			customerType.Id,
			map[int]interface{}{
				customerType.CustomerPropertyTypes[0].Id: "男",
				customerType.CustomerPropertyTypes[1].Id: 20,
			},
		)
		interactor := NewCustomerCreateInteractor()
		response, err := interactor.Handle(request)
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, response.ValidationErrors, 0)
		assert.Equal(t, "山田"+timestampStr, response.CustomerDto.Name)
		assert.Equal(t, customerType.Id, response.CustomerDto.CustomerTypeId)
		expectedProperties := customer.PropertyDto{
			customerType.CustomerPropertyTypes[0].Id: "男",
			customerType.CustomerPropertyTypes[1].Id: 20,
		}
		assert.Equal(t, expectedProperties, response.CustomerDto.Properties)
	})
}

func preCreateCustomerType() (customer_type.CustomerTypeDto, error) {
	timestampStr := utils.CreateTimestampString()

	// カスタマープロパティ新規登録
	var propertyTypeIds []int
	propertyType1, err := preCreateCustomerPropertyType("性別"+timestampStr, "string")
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	propertyTypeIds = append(propertyTypeIds, propertyType1.Id)
	propertyType2, err := preCreateCustomerPropertyType("年齢"+timestampStr, "numeric")
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	propertyTypeIds = append(propertyTypeIds, propertyType2.Id)

	// カスタマータイプ新規登録
	interactor := create2.NewCustomerTypeCreateInteractor(
		db.NewCustomerTypeRepository(),
		db.NewCustomerPropertyTypeRepository(),
	)
	request := create2.NewCustomerTypeCreateUseCaseRequest("個人経営"+timestampStr, propertyTypeIds)
	response, err := interactor.Handle(request)
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	if len(response.ValidationError) > 0 {
		return customer_type.CustomerTypeDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationError)
	}
	return response.CustomerTypeDto, nil
}

func preCreateCustomerPropertyType(name string, propertyType string) (customer_property_type.CustomerPropertyTypeDto, error) {
	interactor := create.NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())
	request := create.NewCustomerPropertyTypeCreateUseCaseRequest(name, propertyType)
	response, err := interactor.Handle(request)
	if err != nil {
		return customer_property_type.CustomerPropertyTypeDto{}, err
	}
	if len(response.ValidationError) > 0 {
		return customer_property_type.CustomerPropertyTypeDto{}, errors.Errorf("validError %+v", response.ValidationError)
	}
	return response.CustomerPropertyTypeDto, nil
}
