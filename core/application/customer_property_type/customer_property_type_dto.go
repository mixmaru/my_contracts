package customer_property_type

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/pkg/errors"
)

type CustomerPropertyTypeDto struct {
	Id   int
	Name string
	Type string
}

func NewCustomerPropertyTypeDtoFromEntity(entity *customer.CustomerPropertyTypeEntity) (CustomerPropertyTypeDto, error) {
	dto := CustomerPropertyTypeDto{}
	dto.Id = entity.Id()
	dto.Name = entity.Name()
	switch entity.PropertyType() {
	case customer.PROPERTY_TYPE_STRING:
		dto.Type = "string"
	case customer.PROPERTY_TYPE_NUMERIC:
		dto.Type = "numeric"
	default:
		return CustomerPropertyTypeDto{}, errors.Errorf("entity.PropertyType()が想定外。entity: %v", entity)
	}

	return dto, nil
}
