package entities

import (
	"database/sql"
	"sort"
	"time"
)

type BillAggregation struct {
	BaseEntity
	billingDate        time.Time
	paymentConfirmedAt sql.NullTime
	billDetails        []*BillDetailEntity
}

func NewBillingAggregation(billingDate time.Time) *BillAggregation {
	return &BillAggregation{
		billingDate: billingDate,
	}
}

func (b *BillAggregation) BillingDate() time.Time {
	return b.billingDate
}

func (b *BillAggregation) PaymentConfirmedAt() (confirmedAt time.Time, isNull bool) {
	return time.Time{}, false
}

func (b *BillAggregation) AddBillDetail(billDetailEntity *BillDetailEntity) {
	// 追加する
	b.billDetails = append(b.billDetails, billDetailEntity)

	// ソートする
	sort.Slice(b.billDetails, func(i, j int) bool {
		return b.billDetails[i].OrderNum() < b.billDetails[j].OrderNum()
	})
}

func (b *BillAggregation) BillDetails() []*BillDetailEntity {
	return nil
}
