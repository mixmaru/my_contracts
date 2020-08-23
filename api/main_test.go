package main

import (
	"encoding/json"
	"fmt"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/application_service/data_transfer_objects"
	"github.com/mixmaru/my_contracts/internal/domains/contracts/repositories/db_connection"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
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
			var registeredUser data_transfer_objects.UserIndividualDto
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
			var registeredUser data_transfer_objects.UserCorporationDto
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
	userAppService := application_service.NewUserApplicationService()
	// 個人ユーザー登録
	savedIndividualUser, validErrors, err := userAppService.RegisterUserIndividual("個人たたろう")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)

	// 法人ユーザー登録
	savedCorporationUser, validErrors, err := userAppService.RegisterUserCorporation("イケイケ株式会社", "法人担当者", "社長次郎")
	assert.NoError(t, err)
	assert.Len(t, validErrors, 0)

	router := newRouter()

	t.Run("個人ユーザー取得", func(t *testing.T) {
		t.Run("userIdをurlで受け取ってそのUser情報を返す", func(t *testing.T) {
			////// 実行
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%v", savedIndividualUser.Id), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 検証
			var loadedUser data_transfer_objects.UserIndividualDto
			// jsonパース
			err = json.Unmarshal(rec.Body.Bytes(), &loadedUser)
			assert.NoError(t, err)
			assert.Equal(t, savedIndividualUser.Id, loadedUser.Id)
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
			req := httptest.NewRequest("GET", fmt.Sprintf("/users/%v", savedCorporationUser.Id), nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			// 検証
			var loadedUser data_transfer_objects.UserCorporationDto
			// jsonパース
			err = json.Unmarshal(rec.Body.Bytes(), &loadedUser)
			assert.NoError(t, err)
			assert.Equal(t, savedCorporationUser.Id, loadedUser.Id)
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
func TestMain_saveContract(t *testing.T) {
	// 重複しない商品名でテストを行う
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 10)
	name := "商品" + suffix

	// 商品登録
	productApp := application_service.NewProductApplicationService()
	productDto, validErrs, err := productApp.Register(name, "200")
	assert.NoError(t, err)
	assert.Len(t, validErrs, 0)
	// ユーザー登録
	userApp := application_service.NewUserApplicationService()
	userDto, validErrs, err := userApp.RegisterUserIndividual("太郎くん")
	assert.NoError(t, err)
	assert.Len(t, validErrs, 0)

	router := newRouter()

	t.Run("顧客IDと商品IDを渡すと契約が作成されて作成された契約データ内容が返る。内部では使用権データも作成されている", func(t *testing.T) {
		// 準備
		body := url.Values{}
		body.Set("user_id", strconv.Itoa(userDto.Id))
		body.Set("product_id", strconv.Itoa(productDto.Id))

		// リクエスト実行
		req := httptest.NewRequest("POST", "/contracts/", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)
		// jsonパース
		var registeredContract data_transfer_objects.ContractDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredContract)
		assert.NoError(t, err)
		assert.Equal(t, userDto.Id, registeredContract.UserId)
		assert.Equal(t, productDto.Id, registeredContract.ProductId)
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
	// 重複しない商品名でテストを行う
	unixNano := time.Now().UnixNano()
	suffix := strconv.FormatInt(unixNano, 10)
	name := "商品" + suffix

	// 検証用データ(商品)登録
	productAppService := application_service.NewProductApplicationService()
	product, validErrs, err := productAppService.Register(name, "100")
	assert.NoError(t, err)
	assert.Len(t, validErrs, 0)

	// 検証用データ(user)登録
	userAppService := application_service.NewUserApplicationService()
	user, validErrs, err := userAppService.RegisterUserCorporation("イケイケ池株式会社", "契約取得用顧客担当", "契約取得用社長")
	assert.NoError(t, err)
	assert.Len(t, validErrs, 0)

	// 検証用データ(契約)登録
	contractAppService := application_service.NewContractApplicationService()
	contract, validErrs, err := contractAppService.Register(user.Id, product.Id, time.Now())
	assert.NoError(t, err)
	assert.Len(t, validErrs, 0)

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

		assert.Equal(t, product.Id, gotContractData.Product.Id)
		assert.Equal(t, name, gotContractData.Product.Name)
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
