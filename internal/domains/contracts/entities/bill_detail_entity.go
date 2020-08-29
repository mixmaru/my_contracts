package entities

type BillDetailEntity struct {
	BaseEntity
	orderNum     int
	rightToUseId int
}

func NewBillingDetailEntity(orderNum, rightToUseId int) *BillDetailEntity {
	return &BillDetailEntity{
		orderNum:     orderNum,
		rightToUseId: rightToUseId,
	}
}

func (b *BillDetailEntity) OrderNum() int {
	return b.orderNum
}

func (b *BillDetailEntity) RightToUseId() int {
	return b.rightToUseId
}
