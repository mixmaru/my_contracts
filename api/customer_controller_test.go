package main

import (
	"bytes"
	"encoding/json"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCustomerController_Create(t *testing.T) {
	router := newRouter()

	t.Run("名前とプロパティIDとプロパティ辞書を渡すと新カスタマーが登録される", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		////// 準備
		postBody := map[string]interface{}{
			"name":             "名前" + timestampstr,
			"customer_type_id": 1,
			"properties": map[string]interface{}{
				"性別": "男",
				"年齢": 20,
			},
		}
		body, _ := json.Marshal(postBody)
		req := httptest.NewRequest("POST", "/customer/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		////// リクエスト実行
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		//// jsonパース
		// todo: check CustomerDto
		//var registeredCustomer CustomerDto
		//err := json.Unmarshal(rec.Body.Bytes(), &registeredCustomer)
		//assert.NoError(t, err)
		//assert.
	})
	//
	//	t.Run("idを渡すとデータが取得できる", func(t *testing.T) {
	//		////// 準備
	//		body := url.Values{}
	//		id := strconv.Itoa(registeredCustomerType.Id)
	//		////// 実行
	//		req := httptest.NewRequest("GET", "/customer_types/"+id, strings.NewReader(body.Encode()))
	//		rec := httptest.NewRecorder()
	//		router.ServeHTTP(rec, req)
	//		////// 検証
	//		assert.Equal(t, http.StatusOK, rec.Code)
	//		// jsonパース
	//		var loadedCustomerType customer_type.CustomerTypeDto
	//		err = json.Unmarshal(rec.Body.Bytes(), &loadedCustomerType)
	//		assert.NoError(t, err)
	//		assert.Equal(t, registeredCustomerType.Id, loadedCustomerType.Id)
	//		assert.Equal(t, "法人"+timestampstr, loadedCustomerType.Name)
	//		assert.Equal(t, preCreateCustomerProperties, loadedCustomerType.CustomerPropertyTypes)
	//	})
	//
	//	t.Run("idに存在しないidを渡すとnot found", func(t *testing.T) {
	//		////// 準備
	//		body := url.Values{}
	//		id := "-10000"
	//		////// 実行
	//		req := httptest.NewRequest("GET", "/customer_types/"+id, strings.NewReader(body.Encode()))
	//		rec := httptest.NewRecorder()
	//		router.ServeHTTP(rec, req)
	//		////// 検証
	//		assert.Equal(t, http.StatusNotFound, rec.Code)
	//	})
	//
	//	t.Run("idに数値以外を渡すとnot found", func(t *testing.T) {
	//		////// 準備
	//		body := url.Values{}
	//		id := "aaa"
	//		////// 実行
	//		req := httptest.NewRequest("GET", "/customer_types/"+id, strings.NewReader(body.Encode()))
	//		rec := httptest.NewRecorder()
	//		router.ServeHTTP(rec, req)
	//		////// 検証
	//		assert.Equal(t, http.StatusNotFound, rec.Code)
	//	})
	//
	//	t.Run("既に存在する名前を登録しようとするとバリデーションエラー。", func(t *testing.T) {
	//		////// 準備
	//		body := url.Values{}
	//		body.Set("name", "法人"+timestampstr) // 上部で既に登録すみ
	//		body.Set("customer_property_ids", strconv.Itoa(preCreateCustomerProperties[0].Id))
	//		body.Add("customer_property_ids", strconv.Itoa(preCreateCustomerProperties[1].Id))
	//
	//		////// 実行
	//		req := httptest.NewRequest("POST", "/customer_types/", strings.NewReader(body.Encode()))
	//		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	//		rec := httptest.NewRecorder()
	//		router.ServeHTTP(rec, req)
	//
	//		////// 検証
	//		assert.Equal(t, http.StatusBadRequest, rec.Code)
	//		// jsonパース
	//		var validErrors map[string][]string
	//
	//		err = json.Unmarshal(rec.Body.Bytes(), &validErrors)
	//		assert.NoError(t, err)
	//		expected := map[string][]string{
	//			"name": []string{
	//				"既に存在する名前です",
	//			},
	//		}
	//		assert.Equal(t, expected, validErrors)
	//	})
	//})
	//
	//t.Run("名前とプロパティIDを渡さないとバリデーションエラー", func(t *testing.T) {
	//	////// 実行
	//	body := url.Values{}
	//	req := httptest.NewRequest("POST", "/customer_types/", strings.NewReader(body.Encode()))
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	//	rec := httptest.NewRecorder()
	//	router.ServeHTTP(rec, req)
	//
	//	////// 検証
	//	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//	// jsonパース
	//	var validErrors map[string][]string
	//
	//	err = json.Unmarshal(rec.Body.Bytes(), &validErrors)
	//	assert.NoError(t, err)
	//	expected := map[string][]string{
	//		"name": []string{
	//			"入力されていません",
	//		},
	//		"customer_property_ids": []string{
	//			"入力されていません",
	//		},
	//	}
	//	assert.Equal(t, expected, validErrors)
	//
	//})
	//
	//t.Run("プロパティIDに数値以外を入れるとバリデーションエラー", func(t *testing.T) {
	//	////// 準備
	//	timestampstr := utils.CreateTimestampString()
	//	body := url.Values{}
	//	body.Set("name", "プロパティIDに変な値"+timestampstr) // 上部で既に登録すみ
	//	body.Set("customer_property_ids", "1a")
	//
	//	////// 実行
	//	req := httptest.NewRequest("POST", "/customer_types/", strings.NewReader(body.Encode()))
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	//	rec := httptest.NewRecorder()
	//	router.ServeHTTP(rec, req)
	//
	//	////// 検証
	//	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//	// jsonパース
	//	var validErrors map[string][]string
	//
	//	err = json.Unmarshal(rec.Body.Bytes(), &validErrors)
	//	assert.NoError(t, err)
	//	expected := map[string][]string{
	//		"customer_property_ids": []string{
	//			"数値ではありません",
	//		},
	//	}
	//	assert.Equal(t, expected, validErrors)
	//})
	//
	//t.Run("プロパティIDに存在しないIDを入れるとバリデーションエラー", func(t *testing.T) {
	//	////// 準備
	//	timestampstr := utils.CreateTimestampString()
	//	body := url.Values{}
	//	body.Set("name", "プロパティIDに存在しないID"+timestampstr) // 上部で既に登録すみ
	//	body.Set("customer_property_ids", "-10000")
	//
	//	////// 実行
	//	req := httptest.NewRequest("POST", "/customer_types/", strings.NewReader(body.Encode()))
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	//	rec := httptest.NewRecorder()
	//	router.ServeHTTP(rec, req)
	//
	//	////// 検証
	//	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//	// jsonパース
	//	var validErrors map[string][]string
	//
	//	err = json.Unmarshal(rec.Body.Bytes(), &validErrors)
	//	assert.NoError(t, err)
	//	expected := map[string][]string{
	//		"customer_property_ids": []string{
	//			"-10000 は存在しないidです",
	//		},
	//	}
	//	assert.Equal(t, expected, validErrors)
	//})
}
