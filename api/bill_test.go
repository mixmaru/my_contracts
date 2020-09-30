package main

import (
	"encoding/json"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createTestDate(t *testing.T) (data_transfer_objects.UserIndividualDto, data_transfer_objects.ProductDto, data_transfer_objects.ContractDto) {
	// user作成
	userApp := application_service.NewUserApplicationService()
	user, validErrors, err := userApp.RegisterUserIndividual("請求実行バッチapiテスト用顧客")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)
	// 商品作成
	productApp := application_service.NewProductApplicationService()
	product, validErrors, err := productApp.Register("商品", "10000")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)
	// 契約作成（使用権も内部で作成されている）
	contractApp := application_service.NewContractApplicationService()
	contract, validErrors, err := contractApp.Register(user.Id, product.Id, utils.CreateJstTime(2020, 6, 1, 12, 30, 26, 111111))
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)
	return user, product, contract
}

func TestMain_executeBilling(t *testing.T) {
	router := newRouter()
	t.Run("指定した日付を基準日にして請求実行を行い作成されたbillデータを返却する", func(t *testing.T) {
		// 事前に請求実行してきれいにしておく
		req := httptest.NewRequest("POST", "/batches/bills/billing?date=20200602", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		////// 準備（商品、ユーザー、契約、使用権を作成する）
		user, _, _ := createTestDate(t)

		// リクエスト実行
		req = httptest.NewRequest("POST", "/batches/bills/billing?date=20200602", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var registeredBills []*data_transfer_objects.BillDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredBills)
		assert.NoError(t, err)

		expectBills := []data_transfer_objects.BillDto{
			{
				BillingDate:             utils.CreateJstTime(2020, 6, 2, 0, 0, 0, 0),
				UserId:                  user.Id,
				PaymentConfirmed:        false,
				PaymentConfirmedAt:      time.Time{},
				TotalAmountExcludingTax: "10000",
				BillDetails: []data_transfer_objects.BillDetailDto{
					data_transfer_objects.BillDetailDto{
						BillingAmount: "10000",
					},
				},
			},
		}
		assert.Len(t, registeredBills, 1)
		assert.NotZero(t, registeredBills[0].Id)
		assert.NotZero(t, registeredBills[0].CreatedAt)
		assert.NotZero(t, registeredBills[0].UpdatedAt)
		assert.Equal(t, expectBills[0].UserId, registeredBills[0].UserId)
		assert.Equal(t, expectBills[0].PaymentConfirmed, registeredBills[0].PaymentConfirmed)
		assert.True(t, registeredBills[0].PaymentConfirmedAt.Equal(expectBills[0].PaymentConfirmedAt))
		assert.Equal(t, expectBills[0].TotalAmountExcludingTax, registeredBills[0].TotalAmountExcludingTax)
		// details
		actualDetails := registeredBills[0].BillDetails
		expectDetails := expectBills[0].BillDetails
		assert.Len(t, actualDetails, 1)
		assert.NotZero(t, actualDetails[0].Id)
		assert.NotZero(t, actualDetails[0].CreatedAt)
		assert.NotZero(t, actualDetails[0].UpdatedAt)
		assert.Equal(t, expectDetails[0].BillingAmount, actualDetails[0].BillingAmount)
	})

	t.Run("指定日付がなければ当日指定で請求実行を行い作成されたbillデータを返却する", func(t *testing.T) {
		// 事前に請求実行してきれいにしておく
		req := httptest.NewRequest("POST", "/batches/bills/billing", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		////// 準備（商品、ユーザー、契約、使用権を作成する）
		user, _, _ := createTestDate(t)

		// リクエスト実行（日付指定なし。上記で新たに作成された使用権の請求データが作成されるはず）
		req = httptest.NewRequest("POST", "/batches/bills/billing", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var registeredBills []*data_transfer_objects.BillDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredBills)
		assert.NoError(t, err)

		expectBills := []data_transfer_objects.BillDto{
			{
				BillingDate:             utils.CreateJstTime(2020, 6, 2, 0, 0, 0, 0),
				UserId:                  user.Id,
				PaymentConfirmed:        false,
				PaymentConfirmedAt:      time.Time{},
				TotalAmountExcludingTax: "10000",
				BillDetails: []data_transfer_objects.BillDetailDto{
					data_transfer_objects.BillDetailDto{
						BillingAmount: "10000",
					},
				},
			},
		}
		assert.Len(t, registeredBills, 1)
		assert.NotZero(t, registeredBills[0].Id)
		assert.NotZero(t, registeredBills[0].CreatedAt)
		assert.NotZero(t, registeredBills[0].UpdatedAt)
		assert.Equal(t, expectBills[0].UserId, registeredBills[0].UserId)
		assert.Equal(t, expectBills[0].PaymentConfirmed, registeredBills[0].PaymentConfirmed)
		assert.True(t, registeredBills[0].PaymentConfirmedAt.Equal(expectBills[0].PaymentConfirmedAt))
		assert.Equal(t, expectBills[0].TotalAmountExcludingTax, registeredBills[0].TotalAmountExcludingTax)
		// details
		actualDetails := registeredBills[0].BillDetails
		expectDetails := expectBills[0].BillDetails
		assert.Len(t, actualDetails, 1)
		assert.NotZero(t, actualDetails[0].Id)
		assert.NotZero(t, actualDetails[0].CreatedAt)
		assert.NotZero(t, actualDetails[0].UpdatedAt)
		assert.Equal(t, expectDetails[0].BillingAmount, actualDetails[0].BillingAmount)
	})

	t.Run("作成されたbillデータがなければ（対象請求がなければ）空配列が返る", func(t *testing.T) {
		// 事前に請求実行してきれいにしておく
		req := httptest.NewRequest("POST", "/batches/bills/billing?date=10010101", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		// リクエスト実行（日付をめっちゃ過去にして実行 => 請求が発生しないはず）
		req = httptest.NewRequest("POST", "/batches/bills/billing?date=10010101", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var registeredBills []*data_transfer_objects.BillDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredBills)
		assert.NoError(t, err)
		assert.Len(t, registeredBills, 0)
	})

	t.Run("指定日付のフォーマットがYYYYMMDDでなければエラーになる", func(t *testing.T) {
		// リクエスト実行（日付指定をaaaaa）で実行
		req := httptest.NewRequest("POST", "/batches/bills/billing?date=aaaaa", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		// jsonパース
		var returnData map[string][]string
		err := json.Unmarshal(rec.Body.Bytes(), &returnData)
		assert.NoError(t, err)
		assert.Len(t, returnData["date"], 1)
		assert.Equal(t, "YYYYMMDDの形式ではありません", returnData["date"][0])
	})
}
