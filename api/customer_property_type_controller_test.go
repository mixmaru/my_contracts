package main

import (
	"encoding/json"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

		t.Run("登録されたカスタマープロパティタイプを取得する", func(t *testing.T) {
			t.Run("idを渡すと取得できる", func(t *testing.T) {
				// 準備
				body := url.Values{}

				// リクエスト実行
				req := httptest.NewRequest("GET", "/customer_property_types/"+strconv.Itoa(registeredCustomerPropertyType.Id), strings.NewReader(body.Encode()))
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				// 検証
				assert.Equal(t, http.StatusOK, rec.Code)
				// jsonパース
				var registeredCustomerPropertyType customer_property_type.CustomerPropertyTypeDto
				err := json.Unmarshal(rec.Body.Bytes(), &registeredCustomerPropertyType)
				assert.NoError(t, err)
				assert.NotZero(t, registeredCustomerPropertyType.Id)
				assert.Equal(t, "好きなアイドル"+timestampstr, registeredCustomerPropertyType.Name)
				assert.Equal(t, "string", registeredCustomerPropertyType.Type)
			})

			t.Run("存在しないidを渡すと404 notfoundになる", func(t *testing.T) {
				// 準備
				body := url.Values{}

				// リクエスト実行
				req := httptest.NewRequest("GET", "/customer_property_types/-2000", strings.NewReader(body.Encode()))
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				// 検証
				assert.Equal(t, http.StatusNotFound, rec.Code)
			})

			t.Run("idに数値以外の文字列を渡すとバリデーションエラー", func(t *testing.T) {
				// 準備
				body := url.Values{}

				// リクエスト実行
				req := httptest.NewRequest("GET", "/customer_property_types/20a", strings.NewReader(body.Encode()))
				rec := httptest.NewRecorder()
				router.ServeHTTP(rec, req)

				// 検証
				assert.Equal(t, http.StatusBadRequest, rec.Code)
				// jsonパース
				var validMessages map[string][]string
				err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
				assert.NoError(t, err)
				expected := map[string][]string{
					"id": {
						"数値ではありません",
					},
				}
				assert.Equal(t, expected, validMessages)
			})
		})
	})
}

func TestCustomerPropertyTypeController_GetAll(t *testing.T) {
	router := newRouter()
	t.Run("全カスタマープロパティが取得できる", func(t *testing.T) {
		t.Run("データがあれば全カスタマープロパティが取得できる", func(t *testing.T) {
			////// 準備
			// 既存データ削除
			conn, err := db.GetConnection()
			assert.NoError(t, err)
			defer conn.Db.Close()
			_, err = conn.Exec("DELETE FROM customers_customer_properties;")
			assert.NoError(t, err)
			_, err = conn.Exec("DELETE FROM customer_types_customer_properties;")
			assert.NoError(t, err)
			_, err = conn.Exec("DELETE FROM customer_properties;")
			assert.NoError(t, err)
			// 既存データ挿入
			timestampStr := utils.CreateTimestampString()
			propertyDto1, err := preCreateCustomerProperty("性別"+timestampStr, "string")
			assert.NoError(t, err)
			propertyDto2, err := preCreateCustomerProperty("年齢"+timestampStr, "numeric")
			assert.NoError(t, err)
			preCreatedProperties := []customer_property_type.CustomerPropertyTypeDto{
				propertyDto1,
				propertyDto2,
			}
			assert.NoError(t, err)

			////// 実行
			body := url.Values{}
			req := httptest.NewRequest("GET", "/customer_property_types/", strings.NewReader(body.Encode()))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusOK, rec.Code)
			// jsonパース
			var loadedCustomerPropertyTypes []customer_property_type.CustomerPropertyTypeDto
			err = json.Unmarshal(rec.Body.Bytes(), &loadedCustomerPropertyTypes)
			assert.NoError(t, err)
			assert.Equal(t, preCreatedProperties, loadedCustomerPropertyTypes)
		})

		t.Run("データがなければ空配列が返る", func(t *testing.T) {
			////// 準備
			// 既存データ削除
			conn, err := db.GetConnection()
			assert.NoError(t, err)
			defer conn.Db.Close()
			_, err = conn.Exec("DELETE FROM customer_types_customer_properties;")
			assert.NoError(t, err)
			_, err = conn.Exec("DELETE FROM customer_properties;")
			assert.NoError(t, err)

			////// 実行
			body := url.Values{}
			req := httptest.NewRequest("GET", "/customer_property_types/", strings.NewReader(body.Encode()))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "[]\n", rec.Body.String())
		})
	})
}
