package main

import (
	"encoding/json"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCustomerPropertyTypeController_Create(t *testing.T) {
	router := newRouter()

	t.Run("名前とタイプを渡すとデータが作成される", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		// 準備
		body := url.Values{}
		body.Set("name", "好きなアイドル"+timestampstr)
		body.Set("type", "string")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/customer_property_types/", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var registeredCustomerPropertyType customer_property_type.CustomerPropertyTypeDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredCustomerPropertyType)
		assert.NoError(t, err)
		assert.NotZero(t, registeredCustomerPropertyType.Id)
		assert.Equal(t, "好きなアイドル"+timestampstr, registeredCustomerPropertyType.Name)
		assert.Equal(t, "string", registeredCustomerPropertyType.Type)

		t.Run("既に存在する名前を登録しようとするとバリデーションエラー。存在しないtypeを入力するとバリデーションエラー", func(t *testing.T) {
			// 準備
			body := url.Values{}
			body.Set("name", "好きなアイドル"+timestampstr) // 上部で既に登録済
			body.Set("type", "aaaa")                 // 適当な値

			// リクエスト実行
			req := httptest.NewRequest("POST", "/customer_property_types/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 検証
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			// jsonパース
			var validMessages map[string][]string
			err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
			assert.NoError(t, err)
			expected := map[string][]string{
				"name": {
					"既に存在する名前です",
				},
				"type": {
					"stringでもnumericでもありません",
				},
			}
			assert.Equal(t, expected, validMessages)
		})
	})
}
