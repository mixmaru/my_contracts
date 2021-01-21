package test_utils

import (
	"github.com/mixmaru/my_contracts/core/application/customer"
	"github.com/mixmaru/my_contracts/core/application/customer/create"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	create2 "github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	create3 "github.com/mixmaru/my_contracts/core/application/customer_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/pkg/errors"
)

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

type PropertyType string

const PROPERTY_TYPE_STRING PropertyType = "string"
const PROPERTY_TYPE_NUMERIC PropertyType = "numeric"

type PropertyParam struct {
	PropertyTypeName string
	PropertyType     PropertyType
}

func PreCreateCustomerPropertyTypeAndCustomerType(
	customerTypeName string,
	customerPropertyTypes []PropertyParam,
) (customer_type.CustomerTypeDto, error) {
	// カスタマープロパティタイプの登録
	propertyIds := []int{}
	for _, propertyType := range customerPropertyTypes {
		propertyDto, err := PreCreateCustomerPropertyType(propertyType.PropertyTypeName, string(propertyType.PropertyType))
		if err != nil {
			return customer_type.CustomerTypeDto{}, err
		}
		propertyIds = append(propertyIds, propertyDto.Id)
	}
	// カスタマータイプの登録
	customerTypeDto, err := PreCreateCustomerType(customerTypeName, propertyIds)
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	return customerTypeDto, nil
}

func PreCreateCustomer(name string, customerTypeId int, properties map[int]interface{}) (customer.CustomerDto, error) {
	request := create.NewCustomerCreateUseCaseRequest(name, customerTypeId, properties)
	intaractor := create.NewCustomerCreateInteractor(db.NewCustomerRepository())
	response, err := intaractor.Handle(request)
	if err != nil {
		return customer.CustomerDto{}, err
	}
	if len(response.ValidationErrors) > 0 {
		return customer.CustomerDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationErrors)
	}
	return response.CustomerDto, nil
}
