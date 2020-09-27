package main

import (
	"encoding/json"
	"fmt"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/domains/contracts/repositories/db_connection"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMain_saveProduct(t *testing.T) {
	conn, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer conn.Db.Close()

	router := newRouter()

	t.Run("商品名と値段を渡すと商品登録して登録データを返す", func(t *testing.T) {
		//////// 準備
		// 重複しない商品名でテストを行う
		unixNano := time.Now().UnixNano()
		suffix := strconv.FormatInt(unixNano, 10)
		name := "商品" + suffix
		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("name", name)
		body.Set("price", "1000.01")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var registeredProduct data_transfer_objects.ProductDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredProduct)
		assert.NoError(t, err)
		assert.Equal(t, name, registeredProduct.Name)
		assert.Equal(t, "1000.01", registeredProduct.Price)
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		t.Run("空文字_要素にから文字を渡すとエラーメッセージを返す", func(t *testing.T) {
			////// 準備
			body := url.Values{}
			body.Set("contact_person_name", "")
			body.Set("president_name", "")

			////// リクエスト実行
			req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			// jsonパース
			var validMessages map[string][]string
			err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
			assert.NoError(t, err)
			expected := map[string][]string{
				"name": {
					"空です",
				},
				"price": {
					"空です",
				},
			}
			assert.Equal(t, expected, validMessages)
		})

		t.Run("文字多すぎだったり_priceがマイナス値だったりするとエラーメッセージを返す", func(t *testing.T) {
			////// 準備
			body := url.Values{}
			body.Set("name", "000000000011111111112222222222333333333344444444445")
			body.Set("price", "-1000")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			// jsonパース
			var validMessages map[string][]string
			err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
			assert.NoError(t, err)
			expected := map[string][]string{
				"name": {
					"50文字より多いです",
				},
				"price": {
					"マイナス値です",
				},
			}
			assert.Equal(t, expected, validMessages)
		})
	})
}

func TestMain_getProduct(t *testing.T) {
	// 重複しない商品名でテストを行う
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 10)
	name := "商品" + suffix

	// 検証用データ登録
	productAppService := application_service.NewProductApplicationService()
	registeredProduct, validErrors, err := productAppService.Register(name, "1000.001")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)

	router := newRouter()

	t.Run("商品IDを受け取って商品データを返す", func(t *testing.T) {
		////// 実行
		req := httptest.NewRequest("GET", fmt.Sprintf("/products/%v", registeredProduct.Id), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 保存したデータを取得
		var gotProductData data_transfer_objects.ProductDto
		err = json.Unmarshal(rec.Body.Bytes(), &gotProductData)
		assert.NoError(t, err)

		////// 検証
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, registeredProduct.Id, gotProductData.Id)
		assert.Equal(t, registeredProduct.Name, gotProductData.Name)
		assert.Equal(t, registeredProduct.Price, gotProductData.Price)
		assert.True(t, registeredProduct.CreatedAt.Equal(gotProductData.CreatedAt))
		assert.True(t, registeredProduct.UpdatedAt.Equal(gotProductData.UpdatedAt))
	})

	t.Run("指定IDの商品が存在しなかった時はNot Roundになる", func(t *testing.T) {
		////// 実行
		req := httptest.NewRequest("GET", "/products/0", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		var jsonValues map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &jsonValues)
		assert.NoError(t, err)
		expect := map[string]string{
			"message": "Not Found",
		}
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, expect, jsonValues)
	})

	t.Run("IDに変な値を入れられた時はNot Foundになる", func(t *testing.T) {
		////// 実行
		req := httptest.NewRequest("GET", "/products/aa99fdsa", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		var jsonValues map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &jsonValues)
		assert.NoError(t, err)
		expect := map[string]string{
			"message": "Not Found",
		}
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, expect, jsonValues)
	})
}
