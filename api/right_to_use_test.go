package main

import (
	"encoding/json"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain_executeRecur(t *testing.T) {
	router := newRouter()
	t.Run("指定した日付を基準日にして時期使用権データを作成する", func(t *testing.T) {
		////// 準備
		// 事前に更新実行してきれいにしておく
		req := httptest.NewRequest("POST", "/batches/right_to_uses/recur?date=20200629", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		// 今回更新対象になる使用権を作成する（2020, 6, 30が使用権の終了日）
		_, _, contract := createTestDate(t)

		// リクエスト実行
		req = httptest.NewRequest("POST", "/batches/right_to_uses/recur?date=20200629", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var actualContracts []*data_transfer_objects.ContractDto
		err := json.Unmarshal(rec.Body.Bytes(), &actualContracts)
		assert.NoError(t, err)
		// 件数
		assert.Len(t, actualContracts, 1)
		// 内容
		actualContract := actualContracts[0]
		assert.Equal(t, contract.Id, actualContract.Id)
		assert.True(t, actualContract.CreatedAt.Equal(contract.CreatedAt))
		assert.True(t, actualContract.UpdatedAt.Equal(contract.UpdatedAt))
		assert.Equal(t, contract.ProductId, actualContract.ProductId)
		assert.True(t, actualContract.BillingStartDate.Equal(contract.BillingStartDate))
		assert.True(t, actualContract.ContractDate.Equal(contract.ContractDate))
		// 使用権
		actualRightToUses := actualContract.RightToUseDtos
		assert.Equal(t, contract.RightToUseDtos[0].Id, actualRightToUses[0].Id)
		assert.NotZero(t, actualRightToUses[1].Id)
		assert.NotZero(t, actualRightToUses[1].CreatedAt)
		assert.NotZero(t, actualRightToUses[1].UpdatedAt)
		assert.True(t, actualRightToUses[1].ValidFrom.Equal(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
		assert.True(t, actualRightToUses[1].ValidTo.Equal(utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0)))
	})
}
