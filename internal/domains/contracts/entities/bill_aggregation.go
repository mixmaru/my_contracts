package entities

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/internal/lib/decimal"
	"github.com/pkg/errors"
	"sort"
	"time"
)

type BillAggregation struct {
	BaseEntity
	billingDate        time.Time
	userId             int
	paymentConfirmedAt sql.NullTime
	billDetails        []*BillDetailEntity
}

func NewBillingAggregation(billingDate time.Time, userId int) *BillAggregation {
	return &BillAggregation{
		billingDate: billingDate,
		userId:      userId,
	}
}

func NewBillingAggregationWithData(
	id int,
	billingDate time.Time,
	userId int,
	paymentConfirmedAt sql.NullTime,
	billDetails []*BillDetailEntity,
	createdAt time.Time,
	updatedAt time.Time,
) *BillAggregation {
	retBill := &BillAggregation{}
	retBill.id = id
	retBill.billingDate = billingDate
	retBill.userId = userId
	retBill.paymentConfirmedAt = paymentConfirmedAt
	retBill.billDetails = billDetails
	retBill.createdAt = createdAt
	retBill.updatedAt = updatedAt
	return retBill
}

func (b *BillAggregation) BillingDate() time.Time {
	return b.billingDate
}

func (b *BillAggregation) UserId() int {
	return b.userId
}

func (b *BillAggregation) PaymentConfirmedAt() (confirmedAt time.Time, isNull bool, err error) {
	if !b.paymentConfirmedAt.Valid {
		// nullの場合
		return time.Time{}, true, nil
	}

	// nullではない場合
	retTime, err := b.paymentConfirmedAt.Value()
	if err != nil {
		return time.Time{}, false, errors.Wrapf(err, "paymentConfirmedAtの取得に失敗しました。billingAggregation: %+v", b)
	}
	return retTime.(time.Time), false, nil
}

func (b *BillAggregation) SetPaymentConfirmedAt(confirmedAt time.Time) error {
	err := b.paymentConfirmedAt.Scan(confirmedAt)
	if err != nil {
		return errors.Wrapf(err, "confirmedAtのセットに失敗しました。billingAggregation: %+v, confirmedAt: %v", b, confirmedAt)
	}
	return nil
}

func (b *BillAggregation) AddBillDetail(billDetailEntity *BillDetailEntity) error {
	// 同じorderNumが既にあればエラーを返す
	for _, billDetail := range b.billDetails {
		if billDetail.orderNum == billDetailEntity.orderNum {
			return errors.Errorf("orderNumは既に存在してます。billDetailEntity: %+v", billDetailEntity)
		}
	}

	// 追加する
	b.billDetails = append(b.billDetails, billDetailEntity)

	// ソートする
	sort.Slice(b.billDetails, func(i, j int) bool {
		return b.billDetails[i].OrderNum() < b.billDetails[j].OrderNum()
	})

	return nil
}

func (b *BillAggregation) BillDetails() []*BillDetailEntity {
	retDetails := make([]*BillDetailEntity, 0, len(b.billDetails))
	for _, detail := range b.billDetails {
		retDetail := *detail
		retDetails = append(retDetails, &retDetail)
	}
	return retDetails
}

func (b *BillAggregation) TotalAmountExcludingTax() decimal.Decimal {
	retAmount := decimal.NewFromInt(0)
	for _, detail := range b.billDetails {
		retAmount = retAmount.Add(detail.billingAmount)
	}
	return retAmount
}
