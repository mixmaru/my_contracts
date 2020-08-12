package domain_service

import (
	"github.com/golang/mock/gomock"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/interfaces/mock_interfaces"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/entities"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/mixmaru/my_contracts/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBillingCalculatorDomainService_BillingAmount(t *testing.T) {
	// テスト用契約を新規作成する
	contract, err := entities.NewContractEntityWithData(
		1,
		2,
		3,
		utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 16, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
		utils.CreateJstTime(2020, 1, 15, 0, 0, 0, 0),
	)
	assert.NoError(t, err)

	////// Productリポジトリモック作成
	// DBから取得される商品データ
	productEntity, err := entities.NewProductEntityWithData(3, "請求金額テスト商品", "1000", time.Now(), time.Now())
	assert.NoError(t, err)
	// mock作成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repositoryMock := mock_interfaces.NewMockIProductRepository(ctrl)
	// GetByIdをモック
	repositoryMock.EXPECT().
		GetById(
			3,
			gomock.Any(),
		).Return(productEntity, nil).
		Times(1)

	// ドメインサービスインスタンス化
	billingDS := NewBillingCalculatorDomainService(repositoryMock)

	// db接続作成
	db, err := db_connection.GetConnection()
	assert.NoError(t, err)

	t.Run("契約初月", func(t *testing.T) {
		t.Run("課金開始日の翌月同日-1より前の日を渡すと_日割り料金が返る", func(t *testing.T) {
			billingAmount, err := billingDS.BillingAmount(contract, utils.CreateJstTime(2020, 2, 14, 15, 0, 0, 0), db)
			assert.NoError(t, err)
			assert.Equal(t, "1000", billingAmount.String())
		})
		//t.Run("課金開始日の翌月同日-1日を渡すと_まるまる1月分の料金が返る", func(t *testing.T) {
		//	billingAmount := BillingAmount()
		//})
	})
	//t.Run("契約翌月以降", func(t *testing.T) {
	//	t.Run("課金開始日の翌月同日-1日より後_翌翌月同日-1の日より前の日を渡すと_日割り料金が返る", func(t *testing.T) {
	//		billingAmount := BillingAmount()
	//	})
	//	t.Run("課金開始日の翌月同日-1日より後_翌翌月同日-1の日を渡すと_まるまる1月分の料金が返る", func(t *testing.T) {
	//		billingAmount := BillingAmount()
	//	})
	//})
}
