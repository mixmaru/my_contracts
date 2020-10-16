package interfaces

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"gopkg.in/gorp.v2"
	"time"
)

type IContractRepository interface {
	Create(contractEntity *entities.ContractEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (contract *entities.ContractEntity, err error)
	GetBillingTargetByBillingDate(billingDate time.Time, executor gorp.SqlExecutor) ([]*entities.ContractEntity, error)
	GetRecurTargets(executeDate time.Time, executor gorp.SqlExecutor) ([]*entities.ContractEntity, error)
	Update(contractEntity *entities.ContractEntity, executor gorp.SqlExecutor) error
	GetHavingExpiredRightToUseContractIds(baseDate time.Time, executor gorp.SqlExecutor) ([]int, error)
}
