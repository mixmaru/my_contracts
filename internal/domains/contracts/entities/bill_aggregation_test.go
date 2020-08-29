package entities

import (
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
		billAggregation.AddBillDetail(billDetailEntity1)
		billAggregation.AddBillDetail(billDetailEntity0)

		// 検証
		assert.Len(t, billAggregation.billDetails, 2)
		assert.Equal(t, billDetailEntity0, billAggregation.billDetails[0])
		assert.Equal(t, billDetailEntity1, billAggregation.billDetails[1])
	})

	t.Run("追加するBillingDetailEntityのOrderNumが重複するとエラーになる", func(t *testing.T) {
	})

	t.Run("同じBillingDetailEntityを2回追加するとエラーになる", func(t *testing.T) {
	})
}

func TestBillAggregation_BillDetails(t *testing.T) {
	t.Run("BillDetailエンティティスライスを取得できる", func(t *testing.T) {

	})
}

func TestBillAggregation_PaymentConfirmedAt(t *testing.T) {
	t.Run("PaymentConfirmedAtがセットされていればtime.Timeで取得できる", func(t *testing.T) {

	})

	t.Run("PaymentConfirmedAtがセットされてなければIsNullがtrueで返る", func(t *testing.T) {

	})
}
