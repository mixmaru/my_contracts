package application_service

import (
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBillApplicationService_ExecuteBilling(t *testing.T) {
	productApp := NewProductApplicationService()
	userApp := NewUserApplicationService()
	contractApp := NewContractApplicationService()
	billApp := NewBillApplicationService()
	t.Run("渡した日付時点で有効な使用権でまだ請求実行データ（billsテーブル）が作成されていないものに請求データを作成する", func(t *testing.T) {
		////// 準備 2020/6/1 ~ 2020/6/30の使用権を作成する
		// 影響するデータを事前削除しておく
		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		_, err = db.Exec("DELETE FROM bill_details")
		assert.NoError(t, err)
		_, err = db.Exec("DELETE FROM right_to_use")
		assert.NoError(t, err)
		// 商品作成
		product, validErrors, err := productApp.Register(utils.CreateUniqProductNameForTest(), "1234")
		assert.NoError(t, err)
		assert.Len(t, validErrors, 0)
		// user作成
		user, validErrors, err := userApp.RegisterUserIndividual("請求実行テスト太郎")
		assert.NoError(t, err)
		assert.Len(t, validErrors, 0)
		// 契約作成（使用権も自動的に作成される）（課金開始日は6/2からになる。）
		_, validErrors, err = contractApp.Register(user.Id, product.Id, utils.CreateJstTime(2020, 6, 1, 0, 0, 0, 0))
		assert.NoError(t, err)
		assert.Len(t, validErrors, 0)

		// 実行 2020/6/2で請求実行する（課金開始日が6/2なので、その日を指定）
		billDtos, err := billApp.ExecuteBilling(utils.CreateJstTime(2020, 6, 2, 0, 0, 0, 0))
		assert.NoError(t, err)

		// 検証 準備で用意した使用権の使用量が請求データに作成されている
		assert.Len(t, billDtos, 1)
		assert.Equal(t, user.Id, billDtos[0].UserId)
		assert.True(t, billDtos[0].BillingDate.Equal(utils.CreateJstTime(2020, 6, 2, 0, 0, 0, 0)))
		assert.False(t, billDtos[0].PaymentConfirmed)
		assert.Zero(t, billDtos[0].PaymentConfirmedAt)
		billDetails := billDtos[0].BillDetails
		assert.Len(t, billDetails, 1)
		assert.Equal(t, "1234", billDetails[0].BillingAmount) // 6/2 ~ 7/2が使用権の期間になってるので日割りにはならない
	})
}
