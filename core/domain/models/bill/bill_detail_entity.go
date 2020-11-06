package bill

import (
	"github.com/mixmaru/my_contracts/lib/decimal"
	"time"
)

type BillDetailEntity struct {
	id            int
	rightToUseId  int
	billingAmount decimal.Decimal
	createdAt     time.Time
	updatedAt     time.Time
}

func NewBillDetailEntity(rightToUseId int, billingAmount decimal.Decimal) *BillDetailEntity {
	return &BillDetailEntity{
		rightToUseId:  rightToUseId,
		billingAmount: billingAmount,
	}
}

func NewBillDetailsEntityWithData(
	id int,
	rightToUseId int,
	billingAmount decimal.Decimal,
	createdAt time.Time,
	updatedAt time.Time,
) *BillDetailEntity {
	retEntity := &BillDetailEntity{}
	retEntity.id = id
	retEntity.rightToUseId = rightToUseId
	retEntity.billingAmount = billingAmount
	retEntity.createdAt = createdAt
	retEntity.updatedAt = updatedAt
	return retEntity
}

func (b *BillDetailEntity) Id() int {
	return b.id
}

func (b *BillDetailEntity) RightToUseId() int {
	return b.rightToUseId
}

func (b *BillDetailEntity) BillingAmount() decimal.Decimal {
	return b.billingAmount
}

func (b *BillDetailEntity) CreatedAt() time.Time {
	return b.createdAt
}

func (b *BillDetailEntity) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *BillDetailEntity) LoadData(id, rightToUseId int, billingAmount decimal.Decimal, createdAt, updatedAt time.Time) {
	b.id = id
	b.rightToUseId = rightToUseId
	b.billingAmount = billingAmount
	b.createdAt = createdAt
	b.updatedAt = updatedAt
}
