package bill

import (
	"github.com/mixmaru/my_contracts/core/domain/models/bill"
	"gopkg.in/gorp.v2"
)

type IBillRepository interface {
	Create(billEntity *bill.BillEntity, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (billEntity *bill.BillEntity, err error)
	GetByUserId(userId int, executor gorp.SqlExecutor) (billEntities []*bill.BillEntity, err error)
}
