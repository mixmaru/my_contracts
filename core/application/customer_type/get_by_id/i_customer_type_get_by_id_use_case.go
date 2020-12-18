package get_by_id

import "github.com/mixmaru/my_contracts/core/application/customer_type"

type ICustomerTypeGetByIdUseCase interface {
	Handle(request *CustomerTypeGetByIdUseCaseRequest) (*CustomerTypeGetByIdUseCaseResponse, error)
}

type CustomerTypeGetByIdUseCaseRequest struct {
	CustomerTypeId int
}

func NewCustomerTypeGetByIdUseCaseRequest(customerTypeId int) *CustomerTypeGetByIdUseCaseRequest {
	return &CustomerTypeGetByIdUseCaseRequest{CustomerTypeId: customerTypeId}
}

type CustomerTypeGetByIdUseCaseResponse struct {
	CustomerTypeDto customer_type.CustomerTypeDto
	ValidationError map[string][]string
}

func NewCustomerTypeGetByIdUseCaseResponse(customerTypeDto customer_type.CustomerTypeDto, validationError map[string][]string) *CustomerTypeGetByIdUseCaseResponse {
	return &CustomerTypeGetByIdUseCaseResponse{CustomerTypeDto: customerTypeDto, ValidationError: validationError}
}
