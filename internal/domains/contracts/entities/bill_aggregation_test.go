package entities

import (
	"fmt"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBillAggregation_AddBillDetail(t *testing.T) {
	t.Run("BillDetailエンティティを追加できる_順番はorderNum順になる", func(t *testing.T) {
		// 準備
		billAggregation := NewBillingAggregation(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0))
		billDetailEntity0 := NewBillingDetailEntity(1, 1)
		billDetailEntity1 := NewBillingDetailEntity(2, 2)

		// 実行
		err := billAggregation.AddBillDetail(billDetailEntity1)
		assert.NoError(t, err)
		err = billAggregation.AddBillDetail(billDetailEntity0)
		assert.NoError(t, err)

		// 検証
		assert.Len(t, billAggregation.billDetails, 2)
		assert.Equal(t, billDetailEntity0, billAggregation.billDetails[0])
		assert.Equal(t, billDetailEntity1, billAggregation.billDetails[1])
	})

	t.Run("追加するBillingDetailEntityのOrderNumが重複するとエラーになる", func(t *testing.T) {
		// 準備
		billAggregation := NewBillingAggregation(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0))
		billDetailEntity0 := NewBillingDetailEntity(1, 1)
		billDetailEntity1 := NewBillingDetailEntity(2, 2)
		billDetailEntity2 := NewBillingDetailEntity(2, 3)

		// 実行・検証
		err := billAggregation.AddBillDetail(billDetailEntity1)
		assert.NoError(t, err)
		err = billAggregation.AddBillDetail(billDetailEntity0)
		assert.NoError(t, err)
		err = billAggregation.AddBillDetail(billDetailEntity2)
		assert.Error(t, err)
	})

	t.Run("同じBillingDetailEntityを2回追加するとエラーになる_orderNumが重複して結果的に実現できてる", func(t *testing.T) {
		// 準備
		billAggregation := NewBillingAggregation(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0))
		billDetailEntity0 := NewBillingDetailEntity(1, 1)
		billDetailEntity1 := NewBillingDetailEntity(2, 2)

		err := billAggregation.AddBillDetail(billDetailEntity0)
		assert.NoError(t, err)
		err = billAggregation.AddBillDetail(billDetailEntity1)
		assert.NoError(t, err)

		// 実行・検証
		err = billAggregation.AddBillDetail(billDetailEntity1)
		assert.Error(t, err)
	})
}

func TestBillAggregation_BillDetails(t *testing.T) {
	t.Run("BillDetailエンティティスライスを取得できる　ただし別メモリにコピーされたやつ　変更されないために", func(t *testing.T) {
		// 準備
		billAggregation := NewBillingAggregation(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0))
		billDetailEntity0 := NewBillingDetailEntity(1, 1)
		billDetailEntity1 := NewBillingDetailEntity(2, 2)
		err := billAggregation.AddBillDetail(billDetailEntity0)
		assert.NoError(t, err)
		err = billAggregation.AddBillDetail(billDetailEntity1)
		assert.NoError(t, err)

		// 実行
		details := billAggregation.BillDetails()

		// 検証
		assert.Len(t, details, 2)
		assert.NotEqual(t, fmt.Sprintf("%p", billAggregation.billDetails[0]), fmt.Sprintf("%p", details[0]))
		assert.EqualValues(t, billAggregation.billDetails[0], details[0])
		assert.NotEqual(t, fmt.Sprintf("%p", billAggregation.billDetails[1]), fmt.Sprintf("%p", details[1]))
		assert.EqualValues(t, billAggregation.billDetails[1], details[1])
	})
}

func TestBillAggregation_PaymentConfirmedAt(t *testing.T) {
	t.Run("PaymentConfirmedAtがセットされていればtime.Timeで取得できる", func(t *testing.T) {
		// 準備
		billAggregation := NewBillingAggregation(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0))
		err := billAggregation.SetPaymentConfirmedAt(utils.CreateJstTime(2020, 1, 15, 15, 0, 0, 0))
		assert.NoError(t, err)

		// 実行
		actual, isNull, err := billAggregation.PaymentConfirmedAt()
		assert.NoError(t, err)

		// 検証
		assert.False(t, isNull)
		assert.True(t, actual.Equal(utils.CreateJstTime(2020, 1, 15, 15, 0, 0, 0)))
	})

	t.Run("PaymentConfirmedAtがセットされてなければIsNullがtrueで返る", func(t *testing.T) {
		// 準備
		billAggregation := NewBillingAggregation(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0))

		// 実行
		actual, isNull, err := billAggregation.PaymentConfirmedAt()
		assert.NoError(t, err)

		// 検証
		assert.True(t, isNull)
		assert.Zero(t, actual)
	})
}
