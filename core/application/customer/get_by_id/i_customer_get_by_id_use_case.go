package get_by_id

import "github.com/mixmaru/my_contracts/core/application/customer"

type ICustomerGetByIdUseCase interface {
	Handle(CustomerGetByIdUseCaseRequest) (CustomerGetByIdResponse, error)
}

type CustomerGetByIdUseCaseRequest struct {
	CustomerId int
}

func NewCustomerGetByIdUseCaseRequest(customerId int) *CustomerGetByIdUseCaseRequest {
	return &CustomerGetByIdUseCaseRequest{CustomerId: customerId}
}

type CustomerGetByIdResponse struct {
	CustomerDto      customer.CustomerDto
	ValidationErrors map[string][]string
}

func NewCustomerGetByIdResponse(customerDto customer.CustomerDto, validationErrors map[string][]string) *CustomerGetByIdResponse {
	return &CustomerGetByIdResponse{CustomerDto: customerDto, ValidationErrors: validationErrors}
}
