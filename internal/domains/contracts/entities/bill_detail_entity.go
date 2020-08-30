package entities

import "github.com/mixmaru/my_contracts/internal/lib/decimal"

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

func (b *BillDetailEntity) OrderNum() int {
	return b.orderNum
}

func (b *BillDetailEntity) RightToUseId() int {
	return b.rightToUseId
}

func (b *BillDetailEntity) BillingAmount() decimal.Decimal {
	return b.billingAmount
}
