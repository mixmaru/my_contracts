package data_mappers

import (
	"github.com/mixmaru/my_contracts/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/pkg/errors"
)

type BillDetailMapper struct {
	Id            int             `db:"id"`
	BillId        int             `db:"bill_id"`
	OrderNum      int             `db:"order_num"`
	RightToUseId  int             `db:"right_to_use_id"`
	BillingAmount decimal.Decimal `db:"billing_amount"`
	CreatedAtUpdatedAtMapper
}

func (b *BillDetailMapper) SetDataToEntity(entity interface{}) error {
	value, ok := entity.(*entities.BillDetailEntity)
	if !ok {
		return errors.Errorf("*entities.BillDetailEntityではないものが渡された。entity: %T", entity)
	}
	value.LoadData(
		b.Id,
		b.RightToUseId,
		b.BillingAmount,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return nil
}
