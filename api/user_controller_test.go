package main

import (
	"encoding/json"
	"fmt"
	"github.com/mixmaru/my_contracts/core/application/users"
	"github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestMain_saveUser(t *testing.T) {
	router := newRouter()
	t.Run("個人顧客", func(t *testing.T) {
		t.Run("typeとnameを受け取って個人ユーザーを登録し_登録した内容を返却する", func(t *testing.T) {
			////// 準備
			// リクエストパラメータ作成
			body := url.Values{}
			body.Set("type", "individual")
			body.Set("name", "個人　太郎")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/users/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusCreated, rec.Code)
			// jsonパース
			var registeredUser users.UserIndividualDto
			err := json.Unmarshal(rec.Body.Bytes(), &registeredUser)
			assert.NoError(t, err)

			assert.NotZero(t, registeredUser.Id)
			assert.Equal(t, "個人　太郎", registeredUser.Name)
			assert.Equal(t, "individual", registeredUser.Type)
			assert.NotZero(t, registeredUser.CreatedAt)
			assert.NotZero(t, registeredUser.UpdatedAt)
		})

		t.Run("バリデーションエラー_nameが空だと登録できずエラーメッセージが返る", func(t *testing.T) {
			////// 準備
			// リクエストパラメータ作成
			body := url.Values{}
			body.Set("type", "individual")
			body.Set("name", "")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/users/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			// jsonパース
			var validMessages map[string][]string
			err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
			assert.NoError(t, err)

			expect := map[string][]string{
				"name": {
					"空です",
				},
			}
			assert.Equal(t, expect, validMessages)
		})
	})

	t.Run("法人顧客", func(t *testing.T) {
		t.Run("typeにcorporationを渡して会社名_担当者名_社長名を渡すと登録されて登録データが返却される", func(t *testing.T) {
			////// 準備
			// リクエストパラメータ作成
			body := url.Values{}
			body.Set("type", "corporation")
			body.Set("corporation_name", "イケイケ株式会社")
			body.Set("contact_person_name", "担当　太郎")
			body.Set("president_name", "社長　太郎")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/users/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusCreated, rec.Code)
			// jsonパース
			var registeredUser users.UserCorporationDto
			err := json.Unmarshal(rec.Body.Bytes(), &registeredUser)
			assert.NoError(t, err)

			assert.NotZero(t, registeredUser.Id)
			assert.Equal(t, "corporation", registeredUser.Type)
			assert.Equal(t, "イケイケ株式会社", registeredUser.CorporationName)
			assert.Equal(t, "担当　太郎", registeredUser.ContactPersonName)
			assert.Equal(t, "社長　太郎", registeredUser.PresidentName)
			assert.NotZero(t, registeredUser.CreatedAt)
			assert.NotZero(t, registeredUser.UpdatedAt)
		})

		t.Run("バリデーションエラー", func(t *testing.T) {
			t.Run("空文字_担当者名や社長名や会社名にから文字を渡すとエラーメッセージが返る", func(t *testing.T) {
				////// 準備
				body := url.Values{}
				body.Set("type", "corporation")
				body.Set("corporation_name", "")
				body.Set("contact_person_name", "")
				body.Set("president_name", "")

				// リクエスト実行
				req := httptest.NewRequest("POST", "/users/", strings.NewReader(body.Encode()))
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
					"corporation_name": {
						"空です",
					},
					"contact_person_name": {
						"空です",
					},
					"president_name": {
						"空です",
					},
				}
				assert.Equal(t, expected, validMessages)
			})

			t.Run("文字多すぎるとエラーメッセージが返る", func(t *testing.T) {
				////// 準備
				body := url.Values{}
				body.Set("type", "corporation")
				body.Set("corporation_name", "000000000011111111112222222222333333333344444444445")
				body.Set("contact_person_name", "000000000011111111112222222222333333333344444444445")
				body.Set("president_name", "００００００００００11111111112222222222333333333344444444445")

				// リクエスト実行
				req := httptest.NewRequest("POST", "/users/", strings.NewReader(body.Encode()))
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
					"corporation_name": {
						"50文字より多いです",
					},
					"contact_person_name": {
						"50文字より多いです",
					},
					"president_name": {
						"50文字より多いです",
					},
				}
				assert.Equal(t, expected, validMessages)
			})
		})
	})

	t.Run("バリデーションエラー_typeに適当な値をいれるとエラーメッセージを返す", func(t *testing.T) {
		////// 準備
		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("type", "aaaa")
		body.Set("name", "個人　太郎")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/users/", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		// jsonパース
		var validErrors map[string][]string
		err := json.Unmarshal(rec.Body.Bytes(), &validErrors)
		assert.NoError(t, err)
		assert.Len(t, validErrors, 1)
		assert.Len(t, validErrors["type"], 1)
		assert.Equal(t, "typeがindividualでもcorporationでもありません。", validErrors["type"][0])
	})
}

func TestMain_getUser(t *testing.T) {

	////// 取得テスト用のデータを登録する
	userRep := db.NewUserRepository()
	// 個人ユーザー登録
	createIndividualIntractor := create.NewUserIndividualCreateInteractor(userRep)
	savedIndividualResponse, err := createIndividualIntractor.Handle(create.NewUserIndividualCreateUseCaseRequest("個人たたろう"))
	assert.NoError(t, err)
	assert.Len(t, savedIndividualResponse.ValidationErrors, 0)

	// 法人ユーザー登録
	createCorprationInteractor := create.NewUserCorporationCreateInteractor(userRep)
	savedCorporationResponse, err := createCorprationInteractor.Handle(create.NewUserCorporationCreateUseCaseRequest("イケイケ株式会社", "法人担当者", "社長次郎"))
	assert.NoError(t, err)
	assert.Len(t, savedCorporationResponse.ValidationErrors, 0)

	router := newRouter()

	t.Run("個人ユーザー取得", func(t *testing.T) {
		t.Run("userIdをurlで受け取ってそのUser情報を返す", func(t *testing.T) {
			////// 実行
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%v", savedIndividualResponse.UserDto.Id), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)

			// 検証
			var loadedUser users.UserIndividualDto
			// jsonパース
			err = json.Unmarshal(rec.Body.Bytes(), &loadedUser)
			assert.NoError(t, err)
			assert.Equal(t, savedIndividualResponse.UserDto.Id, loadedUser.Id)
			assert.Equal(t, "individual", loadedUser.Type)
			assert.Equal(t, "個人たたろう", loadedUser.Name)
			assert.NotZero(t, loadedUser.CreatedAt)
			assert.NotZero(t, loadedUser.UpdatedAt)
		})

		t.Run("指定IDのユーザーが存在しなかったときNotFoundを返す", func(t *testing.T) {
			///// 実行
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%v", -100), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			///// 検証
			assert.Equal(t, http.StatusNotFound, rec.Code)
			// jsonパース
			var jsonValues map[string]string
			err = json.Unmarshal(rec.Body.Bytes(), &jsonValues)
			assert.NoError(t, err)
			expect := map[string]string{
				"message": "Not Found",
			}
			assert.Equal(t, expect, jsonValues)
		})

		t.Run("数値ではない適当な値を入れられたときはNotFoundを返す", func(t *testing.T) {
			////// 実行
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%v", "1a00"), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusNotFound, rec.Code)
			// jsonパース
			var jsonValues map[string]string
			err = json.Unmarshal(rec.Body.Bytes(), &jsonValues)
			assert.NoError(t, err)
			expect := map[string]string{
				"message": "Not Found",
			}
			assert.Equal(t, expect, jsonValues)
		})
	})

	t.Run("法人ユーザー取得", func(t *testing.T) {
		t.Run("UserIdを指定してそのUserの情報を返す", func(t *testing.T) {
			////// 実行
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%v", savedCorporationResponse.UserDto.Id), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 検証
			var loadedUser users.UserCorporationDto
			// jsonパース
			err = json.Unmarshal(rec.Body.Bytes(), &loadedUser)
			assert.NoError(t, err)
			assert.Equal(t, savedCorporationResponse.UserDto.Id, loadedUser.Id)
			assert.Equal(t, "corporation", loadedUser.Type)
			assert.Equal(t, "イケイケ株式会社", loadedUser.CorporationName)
			assert.Equal(t, "法人担当者", loadedUser.ContactPersonName)
			assert.Equal(t, "社長次郎", loadedUser.PresidentName)
			assert.NotZero(t, loadedUser.CreatedAt)
			assert.NotZero(t, loadedUser.UpdatedAt)
		})

		t.Run("指定IDのユーザーが存在しなかったときNotFoundを返す", func(t *testing.T) {
			////// 実行
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%v", -100), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			///// 検証
			assert.Equal(t, http.StatusNotFound, rec.Code)
			// jsonパース
			var jsonValues map[string]string
			err = json.Unmarshal(rec.Body.Bytes(), &jsonValues)
			assert.NoError(t, err)
			expect := map[string]string{
				"message": "Not Found",
			}
			assert.Equal(t, expect, jsonValues)
		})

		t.Run("数値ではない適当な値を入れられたときはNotFoundを返す", func(t *testing.T) {
			////// 実行
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%v", "1a00"), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			////// 検証
			assert.Equal(t, http.StatusNotFound, rec.Code)
			// jsonパース
			var jsonValues map[string]string
			err = json.Unmarshal(rec.Body.Bytes(), &jsonValues)
			assert.NoError(t, err)
			expect := map[string]string{
				"message": "Not Found",
			}
			assert.Equal(t, expect, jsonValues)
		})
	})
}
