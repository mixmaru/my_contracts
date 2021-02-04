package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mixmaru/my_contracts/core/application/customer"
	"github.com/mixmaru/my_contracts/test_utils"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCustomerController_Create(t *testing.T) {
	router := newRouter()
	timestampStr := utils.CreateTimestampString()

	////// 準備
	customerTypeDto, err := test_utils.PreCreateCustomerPropertyTypeAndCustomerType("超お得意様"+timestampStr, []test_utils.PropertyParam{
		{
			PropertyTypeName: "性別" + timestampStr,
			PropertyType:     test_utils.PROPERTY_TYPE_STRING,
		},
		{
			PropertyTypeName: "年齢" + timestampStr,
			PropertyType:     test_utils.PROPERTY_TYPE_NUMERIC,
		},
	})
	assert.NoError(t, err)

	t.Run("名前とプロパティIDとプロパティ辞書を渡すと新カスタマーが登録される", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		////// 準備
		postBody := map[string]interface{}{
			"name":             "名前" + timestampstr,
			"customer_type_id": customerTypeDto.Id,
			"properties": map[int]interface{}{
				customerTypeDto.CustomerPropertyTypes[0].Id: "男",
				customerTypeDto.CustomerPropertyTypes[1].Id: 22,
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
		var registeredCustomer customer.CustomerDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredCustomer)
		assert.NoError(t, err)
		assert.NotZero(t, registeredCustomer.Id)
		assert.Equal(t, "名前"+timestampstr, registeredCustomer.Name)
		assert.Equal(t, customerTypeDto.Id, registeredCustomer.CustomerTypeId)
		expect := customer.PropertyDto{
			customerTypeDto.CustomerPropertyTypes[0].Id: "男",
			customerTypeDto.CustomerPropertyTypes[1].Id: 22.,
		}
		assert.Equal(t, expect, registeredCustomer.Properties)
	})

	t.Run("存在しないカスタマータイプIDを渡すとバリデーションエラーになる", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		////// 準備
		postBody := map[string]interface{}{
			"name":             "名前" + timestampstr,
			"customer_type_id": -100,
			"properties": map[int]interface{}{
				customerTypeDto.CustomerPropertyTypes[0].Id: "男",
				customerTypeDto.CustomerPropertyTypes[1].Id: 22,
			},
		}
		body, _ := json.Marshal(postBody)
		req := httptest.NewRequest("POST", "/customer/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		////// リクエスト実行
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		// jsonパース
		var validationErrors map[string][]string
		err := json.Unmarshal(rec.Body.Bytes(), &validationErrors)
		assert.NoError(t, err)
		expectedValidationErrors := map[string][]string{
			"customer_type_id": []string{"存在しないIDです"},
		}
		assert.Equal(t, expectedValidationErrors, validationErrors)
	})

	t.Run("渡したカスタマータイムに付属するプロパティIDではないIDを渡すとバリデーションエラーになる", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		////// 準備
		// 別のカスタマータイプを登録する
		anotherCustomerTypeDto, err := test_utils.PreCreateCustomerPropertyTypeAndCustomerType(
			"別カスタマータイプ"+timestampStr,
			[]test_utils.PropertyParam{
				{
					PropertyTypeName: "別カスタマータイプ用プロパティ" + timestampStr,
					PropertyType:     test_utils.PROPERTY_TYPE_STRING,
				},
			})
		assert.NoError(t, err)
		// 別のカスタマータイプのプロパティを指定してリクエストしてみる
		postBody := map[string]interface{}{
			"name":             "名前" + timestampstr,
			"customer_type_id": customerTypeDto.Id,
			"properties": map[int]interface{}{
				customerTypeDto.CustomerPropertyTypes[0].Id:        "男",
				customerTypeDto.CustomerPropertyTypes[1].Id:        22,
				anotherCustomerTypeDto.CustomerPropertyTypes[0].Id: "別カスタマープロパティの値",
			},
		}
		body, _ := json.Marshal(postBody)
		req := httptest.NewRequest("POST", "/customer/", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		////// リクエスト実行
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var validationErrors map[string][]string
		err = json.Unmarshal(rec.Body.Bytes(), &validationErrors)
		assert.NoError(t, err)
		expectedValidationErrors := map[string][]string{
			"properties": {
				fmt.Sprintf(
					"id: %v はcustomer_type_id: %v のプロパティタイプではありません",
					anotherCustomerTypeDto.CustomerPropertyTypes[0].Id,
					customerTypeDto.Id,
				)},
		}
		assert.Equal(t, expectedValidationErrors, validationErrors)
	})

	t.Run("存在しないプロパティIDを渡すとバリデーションエラーになる", func(t *testing.T) {

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
	//		assert.Equal(t, preCreateCustomerProperty, loadedCustomerType.CustomerPropertyTypes)
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
	//		body.Set("customer_property_ids", strconv.Itoa(preCreateCustomerProperty[0].Id))
	//		body.Add("customer_property_ids", strconv.Itoa(preCreateCustomerProperty[1].Id))
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

func TestCustomerController_GetById(t *testing.T) {
	router := newRouter()

	////// 準備
	timestampStr := utils.CreateTimestampString()
	customerTypeDto, err := test_utils.PreCreateCustomerPropertyTypeAndCustomerType("超お得意様"+timestampStr, []test_utils.PropertyParam{
		{
			PropertyTypeName: "性別" + timestampStr,
			PropertyType:     test_utils.PROPERTY_TYPE_STRING,
		},
		{
			PropertyTypeName: "年齢" + timestampStr,
			PropertyType:     test_utils.PROPERTY_TYPE_NUMERIC,
		},
	})
	assert.NoError(t, err)

	customerDto, err := test_utils.PreCreateCustomer("顧客名", customerTypeDto.Id, map[int]interface{}{
		customerTypeDto.CustomerPropertyTypes[0].Id: "女",
		customerTypeDto.CustomerPropertyTypes[1].Id: 22.,
	})
	assert.NoError(t, err)

	t.Run("存在するCustomerIdを渡すとCustomerデータが帰ってくる", func(t *testing.T) {
		req := httptest.NewRequest("GET", fmt.Sprintf("/customer/%v/", customerDto.Id), strings.NewReader(""))
		rec := httptest.NewRecorder()

		////// リクエスト実行
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusOK, rec.Code)
		//// jsonパース
		var registeredCustomer customer.CustomerDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredCustomer)
		assert.NoError(t, err)
		assert.NotZero(t, registeredCustomer.Id)
		assert.Equal(t, customerDto.Name, registeredCustomer.Name)
		assert.Equal(t, customerDto.CustomerTypeId, registeredCustomer.CustomerTypeId)
		assert.Equal(t, customerDto.Properties, registeredCustomer.Properties)
	})
}
