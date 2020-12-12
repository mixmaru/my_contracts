package create

import (
	products2 "github.com/mixmaru/my_contracts/core/application/customer_property_type"
)

type ICustomerPropertyTypeCreateUseCase interface {
	Handle(request *CustomerPropertyTypeCreateUseCaseRequest) (*CustomerPropertyTypeCreateUseCaseResponse, error)
}

type CustomerPropertyTypeCreateUseCaseRequest struct {
	Name string
	Type string
}

func NewCustomerPropertyTypeCreateUseCaseRequest(name string, propertyType string) *CustomerPropertyTypeCreateUseCaseRequest {
	//if propertyType != "string" && propertyType != "numeric" {
	//	return nil, errors.Errorf("propertyTypeがstringでもnumericでもありません。propertyType: %v", propertyType)
	//} todo: Handleでバリデーションを実行するからこの時点での確認は不要
	return &CustomerPropertyTypeCreateUseCaseRequest{Name: name, Type: propertyType}
}

type CustomerPropertyTypeCreateUseCaseResponse struct {
	CustomerPropertyTypeDto products2.CustomerPropertyTypeDto
	ValidationError         map[string][]string
}

func NewCustomerPropertyTypeCreateUseCaseResponse(customerPropertyTypeDto products2.CustomerPropertyTypeDto, validationError map[string][]string) *CustomerPropertyTypeCreateUseCaseResponse {
	return &CustomerPropertyTypeCreateUseCaseResponse{CustomerPropertyTypeDto: customerPropertyTypeDto, ValidationError: validationError}
}
