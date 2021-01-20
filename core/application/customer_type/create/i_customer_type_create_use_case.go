package create

import (
	"github.com/mixmaru/my_contracts/core/application/customer_type"
)

type ICustomerTypeCreateUseCase interface {
	Handle(request *CustomerTypeCreateUseCaseRequest) (*CustomerTypeCreateUseCaseResponse, error)
}

type CustomerTypeCreateUseCaseRequest struct {
	Name                    string
	CustomerPropertyTypeIds []int
}

func NewCustomerTypeCreateUseCaseRequest(name string, customerPropertyTypeIds []int) *CustomerTypeCreateUseCaseRequest {
	return &CustomerTypeCreateUseCaseRequest{Name: name, CustomerPropertyTypeIds: customerPropertyTypeIds}
}

type CustomerTypeCreateUseCaseResponse struct {
	CustomerTypeDto  customer_type.CustomerTypeDto
	ValidationErrors map[string][]string
}

func NewCustomerTypeCreateUseCaseResponse(customerTypeDto customer_type.CustomerTypeDto, validationError map[string][]string) *CustomerTypeCreateUseCaseResponse {
	return &CustomerTypeCreateUseCaseResponse{CustomerTypeDto: customerTypeDto, ValidationErrors: validationError}
}
