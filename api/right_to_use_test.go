package main

import (
	"encoding/json"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
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
		// 事前に存在するデータを削除しておく
		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		deleteSql := `
DELETE FROM discount_apply_contract_updates;
DELETE FROM bill_details;
DELETE FROM right_to_use_active;
DELETE FROM right_to_use_history;
DELETE FROM right_to_use;
DELETE FROM contracts;
`
		_, err = db.Exec(deleteSql)
		assert.NoError(t, err)
		// 今回更新対象になる使用権を作成する（2020, 6, 30が使用権の終了日）
		_, _, contract := createTestDate(t)

		// リクエスト実行
		req := httptest.NewRequest("POST", "/batches/right_to_uses/recur?date=20200629", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var actualContracts []*data_transfer_objects.ContractDto
		err = json.Unmarshal(rec.Body.Bytes(), &actualContracts)
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

func TestMain_executeRightToUseArchive(t *testing.T) {
	router := newRouter()
	t.Run("指定した日付を基準日にして期限切れ使用権データをアーカイブする", func(t *testing.T) {
		////// 準備
		// 事前に存在するデータを削除しておく
		db, err := db_connection.GetConnection()
		assert.NoError(t, err)
		deleteSql := `
DELETE FROM discount_apply_contract_updates;
DELETE FROM bill_details;
DELETE FROM right_to_use_active;
DELETE FROM right_to_use_history;
DELETE FROM right_to_use;
DELETE FROM contracts;
`
		_, err = db.Exec(deleteSql)
		assert.NoError(t, err)
		// 今回更新対象になる使用権を作成する（2020, 6, 30が使用権の終了日）
		_, _, contract := createTestDate(t)

		// リクエスト実行
		req := httptest.NewRequest("POST", "/batches/right_to_uses/archive?date=20200730", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var actualRightToUses []*data_transfer_objects.RightToUseDto
		err = json.Unmarshal(rec.Body.Bytes(), &actualRightToUses)
		assert.NoError(t, err)
		// 件数
		assert.Len(t, actualRightToUses, 1)
		// 内容
		actualRightToUse := actualRightToUses[0]
		assert.Equal(t, contract.RightToUseDtos[0].Id, actualRightToUse.Id)
		assert.NotZero(t, actualRightToUse.Id)
		assert.NotZero(t, actualRightToUse.CreatedAt)
		assert.NotZero(t, actualRightToUse.UpdatedAt)
		assert.True(t, actualRightToUse.ValidFrom.Equal(utils.CreateJstTime(2020, 6, 1, 12, 30, 26, 111111000)))
		assert.True(t, actualRightToUse.ValidTo.Equal(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
	})
}
