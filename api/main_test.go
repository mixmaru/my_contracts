package main

import (
	"encoding/json"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMain_saveIndividualUser(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		router := newRouter()

		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("name", "個人　太郎")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/individual_users", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)

		// jsonパース
		var registeredUser data_transfer_objects.UserIndividualDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredUser)
		assert.NoError(t, err)

		assert.Equal(t, "個人　太郎", registeredUser.Name)

	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		router := newRouter()

		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("name", "")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/individual_users", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// jsonパース
		var validMessages map[string][]string
		err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
		assert.NoError(t, err)
	})
}
