package test_utils

import (
	"github.com/mixmaru/my_contracts/core/application/customer"
	"github.com/mixmaru/my_contracts/core/application/customer/create"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	create2 "github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	create3 "github.com/mixmaru/my_contracts/core/application/customer_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
)

func PreCreateCustomer() (customer.CustomerDto, error) {
	timestampStr := utils.CreateTimestampString()
	propertyDto1, err := PreCreateCustomerPropertyType("性別"+timestampStr, "string")
	if err != nil {
		return customer.CustomerDto{}, err
	}
	propertyDto2, err := PreCreateCustomerPropertyType("年齢"+timestampStr, "numeric")
	if err != nil {
		return customer.CustomerDto{}, err
	}
	customerTypeDto, err := PreCreateCustomerType("顧客タイプ"+timestampStr, []int{propertyDto1.Id, propertyDto2.Id})
	if err != nil {
		return customer.CustomerDto{}, err
	}

	request := create.NewCustomerCreateUseCaseRequest(
		"顧客名"+timestampStr,
		customerTypeDto.Id,
		map[int]interface{}{
			propertyDto1.Id: "女",
			propertyDto2.Id: 25.,
		},
	)
	interactor := create.NewCustomerCreateInteractor(db.NewCustomerRepository())
	response, err := interactor.Handle(request)
	if err != nil {
		return customer.CustomerDto{}, err
	}
	if len(response.ValidationErrors) > 0 {
		return customer.CustomerDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationErrors)
	}
	return response.CustomerDto, nil
}

func PreCreateCustomerPropertyType(propertyName string, propertyType string) (customer_property_type.CustomerPropertyTypeDto, error) {
	request := create2.NewCustomerPropertyTypeCreateUseCaseRequest(propertyName, propertyType)
	interactor := create2.NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())
	response, err := interactor.Handle(request)
	if err != nil {
		return customer_property_type.CustomerPropertyTypeDto{}, err
	}
	if len(response.ValidationErrors) > 0 {
		return customer_property_type.CustomerPropertyTypeDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationErrors)
	}
	return response.CustomerPropertyTypeDto, nil
}

func PreCreateCustomerType(customerTypeName string, propertyTypeIds []int) (customer_type.CustomerTypeDto, error) {
	request := create3.NewCustomerTypeCreateUseCaseRequest(customerTypeName, propertyTypeIds)
	interactor := create3.NewCustomerTypeCreateInteractor(db.NewCustomerTypeRepository(), db.NewCustomerPropertyTypeRepository())
	response, err := interactor.Handle(request)
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	if len(response.ValidationErrors) > 0 {
		return customer_type.CustomerTypeDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationErrors)
	}
	return response.CustomerTypeDto, nil
}
