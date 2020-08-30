package entities

import (
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"time"
)

type BillDetailEntity struct {
	BaseEntity
	orderNum      int
	rightToUseId  int
	billingAmount decimal.Decimal
}

func NewBillingDetailEntity(orderNum, rightToUseId int, billingAmount decimal.Decimal) *BillDetailEntity {
	return &BillDetailEntity{
		orderNum:      orderNum,
		rightToUseId:  rightToUseId,
		billingAmount: billingAmount,
	}
}

func NewBillingDetailsEntityWithData(
	id int,
	orderNum int,
	rightToUseId int,
	billingAmount decimal.Decimal,
	createdAt time.Time,
	updatedAt time.Time,
) *BillDetailEntity {
	retEntity := &BillDetailEntity{}
	retEntity.id = id
	retEntity.orderNum = orderNum
	retEntity.rightToUseId = rightToUseId
	retEntity.billingAmount = billingAmount
	retEntity.createdAt = createdAt
	retEntity.updatedAt = updatedAt
	return retEntity
}

func (b *BillDetailEntity) OrderNum() int {
	return b.orderNum
}

func (b *BillDetailEntity) RightToUseId() int {
	return b.rightToUseId
}

func (b *BillDetailEntity) BillingAmount() decimal.Decimal {
	return b.billingAmount
}

func (b *BillDetailEntity) LoadData(id, orderNum, rightToUseId int, billingAmount decimal.Decimal, createdAt, updatedAt time.Time) {
	b.id = id
	b.orderNum = orderNum
	b.rightToUseId = rightToUseId
	b.billingAmount = billingAmount
	b.createdAt = createdAt
	b.updatedAt = updatedAt
}
