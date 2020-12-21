package main

import (
	"encoding/json"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/application/customer_type"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestCustomerTypeController_Create(t *testing.T) {
	router := newRouter()

	t.Run("名前とプロパティIDを渡すとデータが作成される", func(t *testing.T) {
		timestampstr := utils.CreateTimestampString()
		// 準備
		// カスタマープロパティ登録
		preCreateCustomerProperties, err := preCreateCustomerProperties()
		assert.NoError(t, err)

		// カスタマータイプ登録
		body := url.Values{}
		body.Set("name", "法人"+timestampstr)
		body.Set("property_id", strconv.Itoa(preCreateCustomerProperties[0].Id))
		body.Add("property_id", strconv.Itoa(preCreateCustomerProperties[1].Id))

		//// リクエスト実行
		req := httptest.NewRequest("POST", "/customer_types/", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var registeredCustomerType customer_type.CustomerTypeDto
		err = json.Unmarshal(rec.Body.Bytes(), &registeredCustomerType)
		assert.NoError(t, err)
		assert.NotZero(t, registeredCustomerType.Id)
		assert.Equal(t, "法人"+timestampstr, registeredCustomerType.Name)
		assert.Equal(t, preCreateCustomerProperties, registeredCustomerType.CustomerPropertyTypes)

		//t.Run("既に存在する名前を登録しようとするとバリデーションエラー。存在しないtypeを入力するとバリデーションエラー", func(t *testing.T) {
		//	// 準備
		//	body := url.Values{}
		//	body.Set("name", "好きなアイドル"+timestampstr) // 上部で既に登録済
		//	body.Set("type", "aaaa")                 // 適当な値
		//
		//	// リクエスト実行
		//	req := httptest.NewRequest("POST", "/customer_property_types/", strings.NewReader(body.Encode()))
		//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		//	rec := httptest.NewRecorder()
		//	router.ServeHTTP(rec, req)
		//
		//	// 検証
		//	assert.Equal(t, http.StatusBadRequest, rec.Code)
		//	// jsonパース
		//	var validMessages map[string][]string
		//	err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
		//	assert.NoError(t, err)
		//	expected := map[string][]string{
		//		"name": {
		//			"既に存在する名前です",
		//		},
		//		"type": {
		//			"stringでもnumericでもありません",
		//		},
		//	}
		//	assert.Equal(t, expected, validMessages)
		//})
		//
		//t.Run("登録されたカスタマープロパティタイプを取得する", func(t *testing.T) {
		//	t.Run("idを渡すと取得できる", func(t *testing.T) {
		//		// 準備
		//		body := url.Values{}
		//
		//		// リクエスト実行
		//		req := httptest.NewRequest("GET", "/customer_property_types/"+strconv.Itoa(registeredCustomerPropertyType.Id), strings.NewReader(body.Encode()))
		//		rec := httptest.NewRecorder()
		//		router.ServeHTTP(rec, req)
		//
		//		// 検証
		//		assert.Equal(t, http.StatusOK, rec.Code)
		//		// jsonパース
		//		var registeredCustomerPropertyType customer_property_type.CustomerPropertyTypeDto
		//		err := json.Unmarshal(rec.Body.Bytes(), &registeredCustomerPropertyType)
		//		assert.NoError(t, err)
		//		assert.NotZero(t, registeredCustomerPropertyType.Id)
		//		assert.Equal(t, "好きなアイドル"+timestampstr, registeredCustomerPropertyType.Name)
		//		assert.Equal(t, "string", registeredCustomerPropertyType.Type)
		//	})
		//
		//	t.Run("存在しないidを渡すと404 notfoundになる", func(t *testing.T) {
		//		// 準備
		//		body := url.Values{}
		//
		//		// リクエスト実行
		//		req := httptest.NewRequest("GET", "/customer_property_types/-2000", strings.NewReader(body.Encode()))
		//		rec := httptest.NewRecorder()
		//		router.ServeHTTP(rec, req)
		//
		//		// 検証
		//		assert.Equal(t, http.StatusNotFound, rec.Code)
		//	})
		//
		//	t.Run("idに数値以外の文字列を渡すとバリデーションエラー", func(t *testing.T) {
		//		// 準備
		//		body := url.Values{}
		//
		//		// リクエスト実行
		//		req := httptest.NewRequest("GET", "/customer_property_types/20a", strings.NewReader(body.Encode()))
		//		rec := httptest.NewRecorder()
		//		router.ServeHTTP(rec, req)
		//
		//		// 検証
		//		assert.Equal(t, http.StatusBadRequest, rec.Code)
		//		// jsonパース
		//		var validMessages map[string][]string
		//		err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
		//		assert.NoError(t, err)
		//		expected := map[string][]string{
		//			"id": {
		//				"数値ではありません",
		//			},
		//		}
		//		assert.Equal(t, expected, validMessages)
		//	})
		//})
	})
}

func preCreateCustomerProperties() ([]customer_property_type.CustomerPropertyTypeDto, error) {
	timestampStr := utils.CreateTimestampString()
	var retDtos []customer_property_type.CustomerPropertyTypeDto

	for i := 0; i < 2; i++ {
		interactor := create.NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())
		request := create.NewCustomerPropertyTypeCreateUseCaseRequest("性別"+strconv.Itoa(i)+timestampStr, "string")
		response, err := interactor.Handle(request)
		if err != nil {
			return nil, err
		}
		if len(response.ValidationError) > 0 {
			return nil, errors.Errorf("バリデーションエラー。%+v", response.ValidationError)
		}
		retDtos = append(retDtos, response.CustomerPropertyTypeDto)
	}

	return retDtos, nil
}
