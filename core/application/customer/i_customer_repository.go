package customer

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"gopkg.in/gorp.v2"
)

type ICustomerRepository interface {
	Create(customerEntity *customer.CustomerEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (entity *customer.CustomerEntity, err error)
}
