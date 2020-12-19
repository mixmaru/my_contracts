package get_by_ids

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
)

type ICustomerPropertyTypeGetByIdsUseCase interface {
	Handle(request *CustomerPropertyTypeGetByIdsUseCaseRequest) (*CustomerPropertyTypeGetByIdsUseCaseResponse, error)
}

type CustomerPropertyTypeGetByIdsUseCaseRequest struct {
	Ids []int
}

func NewCustomerPropertyTypeGetByIdsUseCaseRequest(ids []int) *CustomerPropertyTypeGetByIdsUseCaseRequest {
	return &CustomerPropertyTypeGetByIdsUseCaseRequest{Ids: ids}
}

type CustomerPropertyTypeGetByIdsUseCaseResponse struct {
	CustomerPropertyTypeDtos []customer_property_type.CustomerPropertyTypeDto
	ValidationError          map[string][]string
}

func NewCustomerPropertyTypeGetByIdsUseCaseResponse(customerPropertyTypeDtos []customer_property_type.CustomerPropertyTypeDto, validationError map[string][]string) *CustomerPropertyTypeGetByIdsUseCaseResponse {
	return &CustomerPropertyTypeGetByIdsUseCaseResponse{CustomerPropertyTypeDtos: customerPropertyTypeDtos, ValidationError: validationError}
}
