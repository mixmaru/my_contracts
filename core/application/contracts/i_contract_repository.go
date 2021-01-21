package contracts

import (
	"github.com/mixmaru/my_contracts/core/domain/models/contract"
	"gopkg.in/gorp.v2"
	"time"
)

type IContractRepository interface {
	Create(contractEntity *contract.ContractEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (contract *contract.ContractEntity, err error)
	GetBillingTargetByBillingDate(billingDate time.Time, executor gorp.SqlExecutor) ([]*contract.ContractEntity, error)
	GetRecurTargets(executeDate time.Time, executor gorp.SqlExecutor) ([]*contract.ContractEntity, error)
	Update(contractEntity *contract.ContractEntity, executor gorp.SqlExecutor) error
	GetHavingExpiredRightToUseContractIds(baseDate time.Time, executor gorp.SqlExecutor) ([]int, error)
}
