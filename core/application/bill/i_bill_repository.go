package bill

import (
	"github.com/mixmaru/my_contracts/core/domain/models/bill"
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"gopkg.in/gorp.v2"
)

type IBillRepository interface {
	Create(billEntity bill.BillEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (billDto *BillDto, err error)
	GetByUserId(userId int, executor gorp.SqlExecutor) (aggregation []*entities.BillAggregation, err error)
}
