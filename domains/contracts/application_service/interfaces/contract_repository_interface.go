package interfaces

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"gopkg.in/gorp.v2"
	"time"
)

type IContractRepository interface {
	Create(contractEntity *entities.ContractEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (contract *entities.ContractEntity, product *entities.ProductEntity, user interface{}, err error)
	GetBillingTargetByBillingDate(billingDate time.Time, executor gorp.SqlExecutor) ([]*entities.ContractEntity, error)
	GetRecurTargets(executeDate time.Time, executor gorp.SqlExecutor) ([]*entities.ContractEntity, error)
}
