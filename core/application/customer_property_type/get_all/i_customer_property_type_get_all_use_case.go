package get_all

import "github.com/mixmaru/my_contracts/core/application/customer_property_type"

type ICustomerPropertyTypeGetAllUseCase interface {
	Handle() (*CustomerPropertyTypeGetAllUseCaseResponse, error)
}

type CustomerPropertyTypeGetAllUseCaseResponse struct {
	CustomerPropertyTypeDtos []customer_property_type.CustomerPropertyTypeDto
}
