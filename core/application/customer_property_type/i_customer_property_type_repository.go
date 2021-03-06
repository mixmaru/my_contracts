package customer_property_type

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"gopkg.in/gorp.v2"
)

type ICustomerPropertyTypeRepository interface {
	Create(entities []*customer.CustomerPropertyTypeEntity, executor gorp.SqlExecutor) (savedIds []int, err error)
	GetAll(executor gorp.SqlExecutor) (propertyTypes []*customer.CustomerPropertyTypeEntity, err error)
	GetByIds(ids []int, executor gorp.SqlExecutor) (propertyTypes []*customer.CustomerPropertyTypeEntity, err error)
	GetByName(name string, executor gorp.SqlExecutor) (propertyType *customer.CustomerPropertyTypeEntity, err error)
}
