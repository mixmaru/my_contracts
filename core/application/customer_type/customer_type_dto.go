package customer_type

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
)

type CustomerTypeDto struct {
	Id                    int
	Name                  string
	CustomerPropertyTypes []customer_property_type.CustomerPropertyTypeDto
}

func NewCustomerTypeDtoFromEntity(
	customerTypeEntity *customer.CustomerTypeEntity,
	customerPropertyTypeEntities []*customer.CustomerPropertyTypeEntity,
) (CustomerTypeDto, error) {
	dto := CustomerTypeDto{}
	dto.Id = customerTypeEntity.Id()
	dto.Name = customerTypeEntity.Name()
	for _, propertyEntity := range customerPropertyTypeEntities {
		propertyDto, err := customer_property_type.NewCustomerPropertyTypeDtoFromEntity(propertyEntity)
		if err != nil {
			return CustomerTypeDto{}, err
		}
		dto.CustomerPropertyTypes = append(dto.CustomerPropertyTypes, propertyDto)
	}
	return dto, nil
}
