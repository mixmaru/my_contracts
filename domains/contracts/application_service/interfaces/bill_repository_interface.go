package interfaces

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"gopkg.in/gorp.v2"
)

type IBillRepository interface {
	Create(billAggregation *entities.BillAggregation, executor gorp.SqlExecutor) (savedId int, err error)
	GetById(id int, executor gorp.SqlExecutor) (aggregation *entities.BillAggregation, err error)
	GetByUserId(userId int, executor gorp.SqlExecutor) (aggregation []*entities.BillAggregation, err error)
}
