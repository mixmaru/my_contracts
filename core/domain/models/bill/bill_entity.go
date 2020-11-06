package bill

import (
	"database/sql"
	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/pkg/errors"
	"time"
)

type BillEntity struct {
	id                 int
	billingDate        time.Time
	userId             int
	paymentConfirmedAt sql.NullTime
	billDetails        []*BillDetailEntity
	createdAt          time.Time
	updatedAt          time.Time
}

func NewBillEntity(billingDate time.Time, userId int) *BillEntity {
	return &BillEntity{
		billingDate: billingDate,
		userId:      userId,
	}
}

func NewBillEntityWithData(
	id int,
	billingDate time.Time,
	userId int,
	paymentConfirmedAt sql.NullTime,
	billDetails []*BillDetailEntity,
	createdAt time.Time,
	updatedAt time.Time,
) *BillEntity {
	retBill := &BillEntity{}
	retBill.id = id
	retBill.billingDate = billingDate
	retBill.userId = userId
	retBill.paymentConfirmedAt = paymentConfirmedAt
	retBill.billDetails = billDetails
	retBill.createdAt = createdAt
	retBill.updatedAt = updatedAt
	return retBill
}

func (b *BillEntity) Id() int {
	return b.id
}

func (b *BillEntity) BillingDate() time.Time {
	return b.billingDate
}

func (b *BillEntity) UserId() int {
	return b.userId
}

func (b *BillEntity) PaymentConfirmedAt() (confirmedAt time.Time, isNull bool, err error) {
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

func (b *BillEntity) SetPaymentConfirmedAt(confirmedAt time.Time) error {
	err := b.paymentConfirmedAt.Scan(confirmedAt)
	if err != nil {
		return errors.Wrapf(err, "confirmedAtのセットに失敗しました。billingAggregation: %+v, confirmedAt: %v", b, confirmedAt)
	}
	return nil
}

func (b *BillEntity) CreatedAt() time.Time {
	return b.createdAt
}

func (b *BillEntity) UpdatedAt() time.Time {
	return b.updatedAt
}

func (b *BillEntity) AddBillDetail(billDetailEntity *BillDetailEntity) error {
	// 同じdetailが既に存在していたらエラー
	for _, detail := range b.billDetails {
		if detail == billDetailEntity {
			return errors.Errorf("既に存在するbillDetailEntityをaddしようとしました。billDetailEntity: %+v", billDetailEntity)
		}
	}

	// 追加する
	b.billDetails = append(b.billDetails, billDetailEntity)
	return nil
}

func (b *BillEntity) BillDetails() []*BillDetailEntity {
	retDetails := make([]*BillDetailEntity, 0, len(b.billDetails))
	for _, detail := range b.billDetails {
		retDetail := *detail
		retDetails = append(retDetails, &retDetail)
	}
	return retDetails
}

func (b *BillEntity) TotalAmountExcludingTax() decimal.Decimal {
	retAmount := decimal.NewFromInt(0)
	for _, detail := range b.billDetails {
		retAmount = retAmount.Add(detail.billingAmount)
	}
	return retAmount
}
