package main

import (
	"encoding/json"
	"fmt"
	"github.com/mixmaru/my_contracts/core/application/contracts"
	create3 "github.com/mixmaru/my_contracts/core/application/contracts/create"
	"github.com/mixmaru/my_contracts/core/application/products/create"
	create2 "github.com/mixmaru/my_contracts/core/application/users/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMain_saveContract(t *testing.T) {
	// 商品登録
	productApp := create.NewProductCreateInteractor(db.NewProductRepository())
	productResponse, err := productApp.Handle(create.NewProductCreateUseCaseRequest("商品", "200"))
	assert.NoError(t, err)
	assert.Len(t, productResponse.ValidationError, 0)
	// ユーザー登録
	userCreateInteractor := create2.NewUserIndividualCreateInteractor(db.NewUserRepository())
	response, err := userCreateInteractor.Handle(create2.NewUserIndividualCreateUseCaseRequest("太郎くん"))
	assert.NoError(t, err)
	assert.Len(t, response.ValidationErrors, 0)
	userDto := response.UserDto

	router := newRouter()

	t.Run("顧客IDと商品IDを渡すと契約が作成されて作成された契約データ内容が返る。内部では使用権データも作成されている", func(t *testing.T) {
		// 準備
		body := url.Values{}
		body.Set("user_id", strconv.Itoa(userDto.Id))
		body.Set("product_id", strconv.Itoa(productResponse.ProductDto.Id))

		// リクエスト実行
		req := httptest.NewRequest("POST", "/contracts/", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var registeredContract contracts.ContractDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredContract)
		assert.NoError(t, err)
		assert.Equal(t, userDto.Id, registeredContract.UserId)
		assert.Equal(t, productResponse.ProductDto.Id, registeredContract.ProductId)
		assert.NotZero(t, registeredContract.Id)
		assert.NotZero(t, registeredContract.ContractDate)
		assert.NotZero(t, registeredContract.BillingStartDate)
		assert.NotZero(t, registeredContract.CreatedAt)
		assert.NotZero(t, registeredContract.UpdatedAt)
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		t.Run("与えられたproduct_idとuser_idが存在しない値だった場合_エラーメッセージが返る", func(t *testing.T) {
			// 準備
			body := url.Values{}
			body.Set("user_id", "-100")
			body.Set("product_id", "-200")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/contracts/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 検証
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			// jsonパース
			var validMessages map[string][]string
			err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
			assert.NoError(t, err)
			expected := map[string][]string{
				"user_id": {
					"存在しません",
				},
				"product_id": {
					"存在しません",
				},
			}
			assert.Equal(t, expected, validMessages)
		})

		t.Run("product_idとuser_idが与えられなかった場合_エラーメッセージが返る", func(t *testing.T) {
			// 準備
			body := url.Values{}

			// リクエスト実行
			req := httptest.NewRequest("POST", "/contracts/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 検証
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			// jsonパース
			var validMessages map[string][]string
			err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
			assert.NoError(t, err)
			expected := map[string][]string{
				"user_id": {
					"数値ではありません",
				},
				"product_id": {
					"数値ではありません",
				},
			}
			assert.Equal(t, expected, validMessages)
		})

		t.Run("product_idとuser_idに数値でないものが与えられた場合_エラーメッセージが返る", func(t *testing.T) {
			// 準備
			body := url.Values{}
			body.Set("user_id", "aaa")
			body.Set("product_id", "-2a00")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/contracts/", strings.NewReader(body.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 検証
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			// jsonパース
			var validMessages map[string][]string
			err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
			assert.NoError(t, err)
			expected := map[string][]string{
				"user_id": {
					"数値ではありません",
				},
				"product_id": {
					"数値ではありません",
				},
			}
			assert.Equal(t, expected, validMessages)
		})
	})
}

func TestMain_getContract(t *testing.T) {
	productRep := db.NewProductRepository()
	userRep := db.NewUserRepository()
	contractRep := db.NewContractRepository()

	// 検証用データ(商品)登録
	productApp := create.NewProductCreateInteractor(productRep)
	productResponse, err := productApp.Handle(create.NewProductCreateUseCaseRequest("商品", "100"))
	assert.NoError(t, err)
	assert.Len(t, productResponse.ValidationError, 0)

	// 検証用データ(user)登録
	userCreateInteractor := create2.NewUserCorporationCreateInteractor(userRep)
	userCreateResponse, err := userCreateInteractor.Handle(create2.NewUserCorporationCreateUseCaseRequest("イケイケ池株式会社", "契約取得用顧客担当", "契約取得用社長"))
	assert.NoError(t, err)
	assert.Len(t, userCreateResponse.ValidationErrors, 0)
	user := userCreateResponse.UserDto

	// 検証用データ(契約)登録
	contractCreateInteractor := create3.NewContractCreateInteractor(userRep, productRep, contractRep)
	contractCreateResponse, err := contractCreateInteractor.Handle(create3.NewContractCreateUseCaseRequest(user.Id, productResponse.ProductDto.Id, time.Now()))
	assert.NoError(t, err)
	assert.Len(t, contractCreateResponse.ValidationErrors, 0)
	contract := contractCreateResponse.ContractDto

	router := newRouter()

	t.Run("GETでcontract_idを渡すと契約情報とユーザー情報が返ってくる", func(t *testing.T) {
		// 実行
		req := httptest.NewRequest("GET", fmt.Sprintf("/contracts/%v", contract.Id), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusOK, rec.Code)
		// 保存したデータを取得
		var gotContractData contractDataForUserCorporation
		err = json.Unmarshal(rec.Body.Bytes(), &gotContractData)
		assert.NoError(t, err)

		assert.NotZero(t, gotContractData.Id)
		assert.NotZero(t, gotContractData.ContractDate)
		assert.NotZero(t, gotContractData.BillingStartDate)
		assert.NotZero(t, gotContractData.CreatedAt)
		assert.NotZero(t, gotContractData.UpdatedAt)

		assert.Equal(t, user.Id, gotContractData.User.Id)
		assert.Equal(t, "corporation", gotContractData.User.Type)
		assert.Equal(t, "イケイケ池株式会社", gotContractData.User.CorporationName)
		assert.Equal(t, "契約取得用顧客担当", gotContractData.User.ContactPersonName)
		assert.Equal(t, "契約取得用社長", gotContractData.User.PresidentName)
		assert.NotZero(t, gotContractData.User.CreatedAt)
		assert.NotZero(t, gotContractData.User.UpdatedAt)

		assert.Equal(t, productResponse.ProductDto.Id, gotContractData.Product.Id)
		assert.Equal(t, "商品", gotContractData.Product.Name)
		assert.Equal(t, "100", gotContractData.Product.Price)
		assert.NotZero(t, gotContractData.Product.CreatedAt)
		assert.NotZero(t, gotContractData.Product.UpdatedAt)
	})

	t.Run("指定IDの契約が存在しなかった時_Not Foundが返る", func(t *testing.T) {
		// 実行
		req := httptest.NewRequest("GET", "/contracts/0", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		var jsonValues map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &jsonValues)
		assert.NoError(t, err)
		expect := map[string]string{
			"message": "Not Found",
		}
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, expect, jsonValues)
	})

	t.Run("IDに変な値を入れられた時_Not Foundが返る", func(t *testing.T) {
		// 実行
		req := httptest.NewRequest("GET", "/contracts/aa99fdsa", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
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

func TestMain_CreateNextRightToUse(t *testing.T) {
	router := newRouter()
	t.Run("指定した日付を基準日にして時期使用権データを作成する", func(t *testing.T) {
		////// 準備
		// 事前に存在するデータを削除しておく
		conn, err := db.GetConnection()
		assert.NoError(t, err)
		deleteSql := `
DELETE FROM discount_apply_contract_updates;
DELETE FROM bill_details;
DELETE FROM right_to_use_active;
DELETE FROM right_to_use_history;
DELETE FROM right_to_use;
DELETE FROM contracts;
`
		_, err = conn.Exec(deleteSql)
		assert.NoError(t, err)
		// 今回更新対象になる使用権を作成する（2020, 6, 30が使用権の終了日）
		_, _, contract := createTestDate(t)

		// リクエスト実行
		req := httptest.NewRequest("POST", "/batches/right_to_uses/recur?date=20200629", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var actualContracts []*contracts.ContractDto
		err = json.Unmarshal(rec.Body.Bytes(), &actualContracts)
		assert.NoError(t, err)
		// 件数
		assert.Len(t, actualContracts, 1)
		// 内容
		actualContract := actualContracts[0]
		assert.Equal(t, contract.Id, actualContract.Id)
		assert.True(t, actualContract.CreatedAt.Equal(contract.CreatedAt))
		assert.True(t, actualContract.UpdatedAt.Equal(contract.UpdatedAt))
		assert.Equal(t, contract.ProductId, actualContract.ProductId)
		assert.True(t, actualContract.BillingStartDate.Equal(contract.BillingStartDate))
		assert.True(t, actualContract.ContractDate.Equal(contract.ContractDate))
		// 使用権
		actualRightToUses := actualContract.RightToUseDtos
		assert.Equal(t, contract.RightToUseDtos[0].Id, actualRightToUses[0].Id)
		assert.NotZero(t, actualRightToUses[1].Id)
		assert.NotZero(t, actualRightToUses[1].CreatedAt)
		assert.NotZero(t, actualRightToUses[1].UpdatedAt)
		assert.True(t, actualRightToUses[1].ValidFrom.Equal(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
		assert.True(t, actualRightToUses[1].ValidTo.Equal(utils.CreateJstTime(2020, 8, 1, 0, 0, 0, 0)))
	})
}

func TestMain_ArchiveExpiredRightToUse(t *testing.T) {
	router := newRouter()
	t.Run("指定した日付を基準日にして期限切れ使用権データをアーカイブする", func(t *testing.T) {
		////// 準備
		// 事前に存在するデータを削除しておく
		conn, err := db.GetConnection()
		assert.NoError(t, err)
		deleteSql := `
DELETE FROM discount_apply_contract_updates;
DELETE FROM bill_details;
DELETE FROM right_to_use_active;
DELETE FROM right_to_use_history;
DELETE FROM right_to_use;
DELETE FROM contracts;
`
		_, err = conn.Exec(deleteSql)
		assert.NoError(t, err)
		// 今回更新対象になる使用権を作成する（2020, 6, 30が使用権の終了日（20200701がValidToの値））
		_, _, contract := createTestDate(t)

		// リクエスト実行
		req := httptest.NewRequest("POST", "/batches/right_to_uses/archive?date=20200701", nil)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		////// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var actualRightToUses []*contracts.RightToUseDto
		err = json.Unmarshal(rec.Body.Bytes(), &actualRightToUses)
		assert.NoError(t, err)
		// 件数
		assert.Len(t, actualRightToUses, 1)
		// 内容
		actualRightToUse := actualRightToUses[0]
		assert.Equal(t, contract.RightToUseDtos[0].Id, actualRightToUse.Id)
		assert.NotZero(t, actualRightToUse.Id)
		assert.NotZero(t, actualRightToUse.CreatedAt)
		assert.NotZero(t, actualRightToUse.UpdatedAt)
		assert.True(t, actualRightToUse.ValidFrom.Equal(utils.CreateJstTime(2020, 6, 1, 12, 30, 26, 111111000)))
		assert.True(t, actualRightToUse.ValidTo.Equal(utils.CreateJstTime(2020, 7, 1, 0, 0, 0, 0)))
	})
}
