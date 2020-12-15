package create

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
)

type ICustomerPropertyTypeGetByIdUseCase interface {
	Handle(request *CustomerPropertyTypeGetByIdUseCaseRequest) (*CustomerPropertyTypeGetByIdUseCaseResponse, error)
}

type CustomerPropertyTypeGetByIdUseCaseRequest struct {
	Id int
}

func NewCustomerPropertyTypeGetByIdUseCaseRequest(id int) *CustomerPropertyTypeGetByIdUseCaseRequest {
	return &CustomerPropertyTypeGetByIdUseCaseRequest{Id: id}
}

type CustomerPropertyTypeGetByIdUseCaseResponse struct {
	CustomerPropertyTypeDto customer_property_type.CustomerPropertyTypeDto
	ValidationError         map[string][]string
}

func NewCustomerPropertyTypeGetByIdUseCaseResponse(customerPropertyTypeDto customer_property_type.CustomerPropertyTypeDto, validationError map[string][]string) *CustomerPropertyTypeGetByIdUseCaseResponse {
	return &CustomerPropertyTypeGetByIdUseCaseResponse{CustomerPropertyTypeDto: customerPropertyTypeDto, ValidationError: validationError}
}
