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

		expect := map[string][]string{
			"name": []string{
				"空です",
			},
		}
		assert.Equal(t, expect, validMessages)
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
		var jsonValues map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &jsonValues)
		assert.NoError(t, err)

		expect := map[string]string{
			"message": "Not Found",
		}
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, expect, jsonValues)
	})

	t.Run("IDに変な値を入れられた時", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", "/individual_users/aa1aa", nil)
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

func TestMain_saveCorporationUser(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		router := newRouter()

		//////// データ登録テスト

		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("contact_person_name", "担当　太郎")
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

	t.Run("バリデーションエラー", func(t *testing.T) {
		router := newRouter()
		// リクエストパラメータ作成

		t.Run("空文字", func(t *testing.T) {
			body := url.Values{}
			body.Set("contact_person_name", "")
			body.Set("president_name", "")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/corporation_users/", strings.NewReader(body.Encode()))
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
				"contact_person_name": []string{
					"空です",
				},
				"president_name": []string{
					"空です",
				},
			}
			assert.Equal(t, expected, validMessages)
		})

		t.Run("文字多すぎ", func(t *testing.T) {
			body := url.Values{}
			body.Set("contact_person_name", "000000000011111111112222222222333333333344444444445")
			body.Set("president_name", "００００００００００11111111112222222222333333333344444444445")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/corporation_users/", strings.NewReader(body.Encode()))
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
				"contact_person_name": []string{
					"50文字より多いです",
				},
				"president_name": []string{
					"50文字より多いです",
				},
			}
			assert.Equal(t, expected, validMessages)
		})
	})
}

func TestMain_getCorporationUser(t *testing.T) {
	// 検証用データ登録
	router := newRouter()
	body := url.Values{}
	body.Set("contact_person_name", "担当　太郎")
	body.Set("president_name", "社長　太郎")

	// リクエスト実行
	req := httptest.NewRequest("POST", "/corporation_users/", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	// 検証
	assert.Equal(t, http.StatusCreated, rec.Code)

	// 保存したデータを取得
	var registeredUser data_transfer_objects.UserCorporationDto
	err := json.Unmarshal(rec.Body.Bytes(), &registeredUser)
	assert.NoError(t, err)

	t.Run("正常系", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", fmt.Sprintf("/corporation_users/%v", registeredUser.Id), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 保存したデータを取得
		var gotUserData data_transfer_objects.UserCorporationDto
		err = json.Unmarshal(rec.Body.Bytes(), &gotUserData)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, registeredUser, gotUserData)
	})

	t.Run("指定IDの個人顧客が存在しなかった時", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", "/corporation_users/0", nil)
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

	t.Run("IDに変な値を入れられた時", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", "/corporation_users/aa99fdsa", nil)
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

func TestMain_saveProduct(t *testing.T) {
	// 同盟商品は登録できないので予め削除
	conn, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer conn.Db.Close()

	_, err = conn.Exec("delete from products where name = 'A商品'")
	assert.NoError(t, err)

	t.Run("正常系", func(t *testing.T) {
		router := newRouter()

		//////// データ登録テスト

		// リクエストパラメータ作成
		body := url.Values{}
		body.Set("name", "A商品")
		body.Set("price", "1000.01")

		// リクエスト実行
		req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 検証
		assert.Equal(t, http.StatusCreated, rec.Code)

		// jsonパース
		var registeredProduct data_transfer_objects.ProductDto
		err := json.Unmarshal(rec.Body.Bytes(), &registeredProduct)
		assert.NoError(t, err)

		assert.Equal(t, "A商品", registeredProduct.Name)
		assert.Equal(t, "1000.01", registeredProduct.Price)
	})

	t.Run("バリデーションエラー", func(t *testing.T) {
		router := newRouter()
		// リクエストパラメータ作成

		t.Run("空文字", func(t *testing.T) {
			body := url.Values{}
			body.Set("contact_person_name", "")
			body.Set("president_name", "")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
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
				"name": []string{
					"空です",
				},
				"price": []string{
					"空です",
				},
			}
			assert.Equal(t, expected, validMessages)
		})

		t.Run("文字多すぎ　priceがマイナス値", func(t *testing.T) {
			body := url.Values{}
			body.Set("name", "000000000011111111112222222222333333333344444444445")
			body.Set("price", "-1000")

			// リクエスト実行
			req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
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
				"name": []string{
					"50文字より多いです",
				},
				"price": []string{
					"マイナス値です",
				},
			}
			assert.Equal(t, expected, validMessages)
		})
	})
}

func TestMain_getProduct(t *testing.T) {
	// 重複商品名は登録できないので予め削除
	conn, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer conn.Db.Close()

	_, err = conn.Exec("delete from products where name = 'A商品'")
	assert.NoError(t, err)

	// 検証用データ登録
	router := newRouter()
	body := url.Values{}
	body.Set("name", "A商品")
	body.Set("price", "1000.001")

	// リクエスト実行
	req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	// 検証
	assert.Equal(t, http.StatusCreated, rec.Code)

	// 保存したデータを取得
	var registeredProduct data_transfer_objects.ProductDto
	err = json.Unmarshal(rec.Body.Bytes(), &registeredProduct)
	assert.NoError(t, err)

	t.Run("正常系", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", fmt.Sprintf("/products/%v", registeredProduct.Id), nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// 保存したデータを取得
		var gotProductData data_transfer_objects.ProductDto
		err = json.Unmarshal(rec.Body.Bytes(), &gotProductData)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, registeredProduct, gotProductData)
	})

	t.Run("指定IDの商品が存在しなかった時", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", "/products/0", nil)
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

	t.Run("IDに変な値を入れられた時", func(t *testing.T) {
		router := newRouter()
		req := httptest.NewRequest("GET", "/products/aa99fdsa", nil)
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
func TestMain_saveContract(t *testing.T) {
	conn, err := db_connection.GetConnection()
	assert.NoError(t, err)
	defer conn.Db.Close()

	// 同盟商品は登録できないので予め契約とともに削除
	_, err = conn.Exec(
		"delete from contracts " +
			"using products " +
			"where products.id = contracts.product_id " +
			"and products.name = 'ab商品' ")
	assert.NoError(t, err)
	_, err = conn.Exec("delete from products where name = 'ab商品'")
	assert.NoError(t, err)

	// 商品登録
	productApp := application_service.NewProductApplicationService()
	productDto, validErrs, err := productApp.Register("ab商品", "200")
	assert.NoError(t, err)
	assert.Len(t, validErrs, 0)
	// ユーザー登録
	userApp := application_service.NewUserApplicationService()
	userDto, validErrs, err := userApp.RegisterUserIndividual("太郎くん")
	assert.NoError(t, err)
	assert.Len(t, validErrs, 0)

	t.Run("正常系", func(t *testing.T) {
		router := newRouter()

		//////// データ登録テスト

		// リクエストパラメータ作成
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
		assert.NotZero(t, registeredContract.CreatedAt)
		assert.NotZero(t, registeredContract.UpdatedAt)
	})

	//t.Run("バリデーションエラー", func(t *testing.T) {
	//	router := newRouter()
	//	// リクエストパラメータ作成
	//
	//	t.Run("空文字", func(t *testing.T) {
	//		body := url.Values{}
	//		body.Set("contact_person_name", "")
	//		body.Set("president_name", "")
	//
	//		// リクエスト実行
	//		req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
	//		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	//		rec := httptest.NewRecorder()
	//		router.ServeHTTP(rec, req)
	//
	//		// 検証
	//		assert.Equal(t, http.StatusBadRequest, rec.Code)
	//
	//		// jsonパース
	//		var validMessages map[string][]string
	//		err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
	//		assert.NoError(t, err)
	//
	//		expected := map[string][]string{
	//			"name": []string{
	//				"空です",
	//			},
	//			"price": []string{
	//				"空です",
	//			},
	//		}
	//		assert.Equal(t, expected, validMessages)
	//	})
	//
	//	t.Run("文字多すぎ　priceがマイナス値", func(t *testing.T) {
	//		body := url.Values{}
	//		body.Set("name", "000000000011111111112222222222333333333344444444445")
	//		body.Set("price", "-1000")
	//
	//		// リクエスト実行
	//		req := httptest.NewRequest("POST", "/products/", strings.NewReader(body.Encode()))
	//		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //formからの入力ということを指定してるっぽい
	//		rec := httptest.NewRecorder()
	//		router.ServeHTTP(rec, req)
	//
	//		// 検証
	//		assert.Equal(t, http.StatusBadRequest, rec.Code)
	//
	//		// jsonパース
	//		var validMessages map[string][]string
	//		err := json.Unmarshal(rec.Body.Bytes(), &validMessages)
	//		assert.NoError(t, err)
	//
	//		expected := map[string][]string{
	//			"name": []string{
	//				"50文字より多いです",
	//			},
	//			"price": []string{
	//				"マイナス値です",
	//			},
	//		}
	//		assert.Equal(t, expected, validMessages)
	//	})
	//})
}
