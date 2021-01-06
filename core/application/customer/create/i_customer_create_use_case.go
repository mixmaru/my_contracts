package create

import "github.com/mixmaru/my_contracts/core/application/customer"

type ICustomerCreateUseCase interface {
	Handle(request *CustomerCreateUseCaseRequest) (*CustomerCreateUseCaseResponse, error)
}

type CustomerCreateUseCaseRequest struct {
	Name           string
	CustomerTypeId int
	Properties     map[string]interface{}
}

func NewCustomerCreateUseCaseRequest(name string, customerTypeId int, properties map[string]interface{}) *CustomerCreateUseCaseRequest {
	return &CustomerCreateUseCaseRequest{Name: name, CustomerTypeId: customerTypeId, Properties: properties}
}

type CustomerCreateUseCaseResponse struct {
	CustomerDto      customer.CustomerDto
	ValidationErrors map[string][]string
}
