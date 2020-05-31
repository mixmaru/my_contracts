package main

import (
	"encoding/json"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestMain_saveIndividualUser_getIndividualUser(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		router := newRouter()

		//////// データ登録テスト

		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("name", "個人　太郎")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/individual_users/", strings.NewReader(body.Encode()))
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

		///////// データ取得テスト
		// リクエスト実行
		registeredId := strconv.Itoa(registeredUser.Id)
		req = httptest.NewRequest("GET", "/individual_users/"+registeredId, nil)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusOK, rec.Code)

		var loadedUser data_transfer_objects.UserIndividualDto
		// jsonパース
		err = json.Unmarshal(rec.Body.Bytes(), &loadedUser)
		assert.NoError(t, err)

		assert.Equal(t, registeredUser, loadedUser)
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		router := newRouter()

		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("name", "")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/individual_users/", strings.NewReader(body.Encode()))
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
func TestMain_getIndividualUser(t *testing.T) {
	// 正常系はTestMain_saveIndividualUser_getIndividualUserでテスト済
	t.Run("指定IDの個人顧客が存在しなかった時", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", "/individual_users/0", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("IDに変な値を入れられた時", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", "/individual_users/aa1aa", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestMain_saveCorporationUser(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		router := newRouter()

		//////// データ登録テスト

		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("contact_name", "担当　太郎")
		body.Set("president_name", "社長　太郎")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/corporation_users/", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)

		// jsonパース
		var registeredUser data_transfer_objects.UserCorporationDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredUser)
		assert.NoError(t, err)

		assert.Equal(t, "担当　太郎", registeredUser.ContactPersonName)
		assert.Equal(t, "社長　太郎", registeredUser.PresidentName)
	})

	//t.Run("バリデーションエラー", func(t *testing.T) {
	//	router := newRouter()
	//
	//	// リクエストパラメータ作成
	//	body := url.Values{}
	//	body.Set("name", "")
	//
	//	// リクエスト実行
	//	req := httptest.NewRequest("POST", "/individual_users/", strings.NewReader(body.Encode()))
	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	//	rec := httptest.NewRecorder()
	//	router.ServeHTTP(rec, req)
	//
	//	// 検証
	//	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//
	//	// jsonパース
	//	var validMessages map[string][]string
	//	err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
	//	assert.NoError(t, err)
	//})
}
