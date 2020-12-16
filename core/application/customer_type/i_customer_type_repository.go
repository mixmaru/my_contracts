package customer_type

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"gopkg.in/gorp.v2"
)

type ICustomerTypeRepository interface {
	Create(customerTypeEntity *customer.CustomerTypeEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (entity *customer.CustomerTypeEntity, err error)
}
