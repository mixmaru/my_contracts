package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mixmaru/my_contracts/core/application/customer"
	create2 "github.com/mixmaru/my_contracts/core/application/customer/create"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	"github.com/mixmaru/my_contracts/core/application/customer_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCustomerController_Create(t *testing.T) {
	router := newRouter()

	////// 準備
	customerTypeDto, err := preCreateCustomerPropertyTypeAndCustomerType()
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

func preCreateCustomerPropertyTypeAndCustomerType() (customer_type.CustomerTypeDto, error) {
	timestampStr := utils.CreateTimestampString()
	// カスタマープロパティタイプの登録
	propertyDto1, err := preCreateCustomerProperty("性別"+timestampStr, "string")
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	propertyDto2, err := preCreateCustomerProperty("年齢"+timestampStr, "numeric")
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	// カスタマータイプの登録
	propertyIds := []int{
		propertyDto1.Id,
		propertyDto2.Id,
	}
	customerTypeDto, err := preCreateCustomerType("超お得意様"+timestampStr, propertyIds)
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	return customerTypeDto, nil
}

func preCreateCustomerType(name string, customerPropertyTypeIds []int) (customer_type.CustomerTypeDto, error) {
	request := create.NewCustomerTypeCreateUseCaseRequest(name, customerPropertyTypeIds)
	interactor := create.NewCustomerTypeCreateInteractor(db.NewCustomerTypeRepository(), db.NewCustomerPropertyTypeRepository())
	response, err := interactor.Handle(request)
	if err != nil {
		return customer_type.CustomerTypeDto{}, err
	}
	if len(response.ValidationError) > 0 {
		return customer_type.CustomerTypeDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationError)
	}
	return response.CustomerTypeDto, nil
}

func TestCustomerController_GetById(t *testing.T) {
	router := newRouter()

	////// 準備
	customerDto, err := preCreateCustomer()
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

func preCreateCustomer() (customer.CustomerDto, error) {
	// CustomerType登録
	customerTypeDto, err := preCreateCustomerPropertyTypeAndCustomerType()
	if err != nil {
		return customer.CustomerDto{}, err
	}
	// Customer登録
	// property作成
	properties := map[int]interface{}{
		customerTypeDto.CustomerPropertyTypes[0].Id: "男",
		customerTypeDto.CustomerPropertyTypes[1].Id: 39,
	}
	request := create2.NewCustomerCreateUseCaseRequest("顧客名", customerTypeDto.Id, properties)
	intaractor := create2.NewCustomerCreateInteractor(db.NewCustomerRepository())
	response, err := intaractor.Handle(request)
	if err != nil {
		return customer.CustomerDto{}, err
	}
	if len(response.ValidationErrors) > 0 {
		return customer.CustomerDto{}, errors.Errorf("バリデーションエラー。%+v", response.ValidationErrors)
	}
	return response.CustomerDto, nil
}
