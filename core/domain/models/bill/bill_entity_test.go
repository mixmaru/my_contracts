package bill

import (
	"fmt"
	"github.com/mixmaru/my_contracts/lib/decimal"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBillEntity_AddBillDetail(t *testing.T) {
	t.Run("BillDetailエンティティを追加できる_順番は追加順になる", func(t *testing.T) {
		// 準備
		billAggregation := NewBillEntity(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), 1)
		billDetailEntity0 := NewBillDetailEntity(1, decimal.NewFromInt(100))
		billDetailEntity1 := NewBillDetailEntity(2, decimal.NewFromInt(100))

		// 実行
		err := billAggregation.AddBillDetail(billDetailEntity0)
		assert.NoError(t, err)
		err = billAggregation.AddBillDetail(billDetailEntity1)
		assert.NoError(t, err)

		// 検証
		assert.Len(t, billAggregation.billDetails, 2)
		assert.Equal(t, billDetailEntity0, billAggregation.billDetails[0])
		assert.Equal(t, billDetailEntity1, billAggregation.billDetails[1])
	})

	t.Run("同じBillingDetailEntityを2回追加するとエラーになる", func(t *testing.T) {
		// 準備
		billAggregation := NewBillEntity(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), 1)
		billDetailEntity0 := NewBillDetailEntity(1, decimal.NewFromInt(100))
		billDetailEntity1 := NewBillDetailEntity(2, decimal.NewFromInt(100))

		err := billAggregation.AddBillDetail(billDetailEntity0)
		assert.NoError(t, err)
		err = billAggregation.AddBillDetail(billDetailEntity1)
		assert.NoError(t, err)

		// 実行・検証
		err = billAggregation.AddBillDetail(billDetailEntity1)
		assert.Error(t, err)
	})
}

func TestBillEntity_BillDetails(t *testing.T) {
	t.Run("BillDetailエンティティスライスを取得できる　ただし別メモリにコピーされたやつ　変更されないために", func(t *testing.T) {
		// 準備
		billAggregation := NewBillEntity(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), 1)
		billDetailEntity0 := NewBillDetailEntity(1, decimal.NewFromInt(100))
		billDetailEntity1 := NewBillDetailEntity(2, decimal.NewFromInt(100))
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

func TestBillEntity_PaymentConfirmedAt(t *testing.T) {
	t.Run("PaymentConfirmedAtがセットされていればtime.Timeで取得できる", func(t *testing.T) {
		// 準備
		billAggregation := NewBillEntity(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), 1)
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
		billAggregation := NewBillEntity(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), 1)

		// 実行
		actual, isNull, err := billAggregation.PaymentConfirmedAt()
		assert.NoError(t, err)

		// 検証
		assert.True(t, isNull)
		assert.Zero(t, actual)
	})
}

func TestBillEntity_TotalAmountExcludingTax(t *testing.T) {
	t.Run("税抜き請求合計金額を取得できる", func(t *testing.T) {
		// 準備
		billAggregation := NewBillEntity(utils.CreateJstTime(2020, 1, 1, 0, 0, 0, 0), 1)
		billDetailEntity0 := NewBillDetailEntity(1, decimal.NewFromInt(100))
		billDetailEntity1 := NewBillDetailEntity(2, decimal.NewFromInt(1000))
		err := billAggregation.AddBillDetail(billDetailEntity0)
		assert.NoError(t, err)
		err = billAggregation.AddBillDetail(billDetailEntity1)
		assert.NoError(t, err)

		// 実行
		actual := billAggregation.TotalAmountExcludingTax()

		// 検証
		assert.Equal(t, "1100", actual.String())
	})
}
