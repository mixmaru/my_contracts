package get_all

import (
	"github.com/mixmaru/my_contracts/core/application/customer_property_type"
	"github.com/mixmaru/my_contracts/core/application/customer_property_type/create"
	"github.com/mixmaru/my_contracts/core/infrastructure/db"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestCustomerPropertyTypeGetAllInteractor_Handle(t *testing.T) {
	t.Run("全てのカスタマープロパティタイプがid順に取得できる", func(t *testing.T) {
		////// 準備
		// 既存データ削除
		conn, err := db.GetConnection()
		assert.NoError(t, err)
		defer conn.Db.Close()
		_, err = conn.Exec("DELETE FROM customers_customer_properties;")
		assert.NoError(t, err)
		_, err = conn.Exec("DELETE FROM customer_types_customer_properties;")
		assert.NoError(t, err)
		_, err = conn.Exec("DELETE FROM customer_properties;")
		assert.NoError(t, err)
		// 既存データ挿入
		preCreatedCustomerPropertyTypes, err := preCreateCustomerPropertyTypes()
		assert.NoError(t, err)

		////// 実行
		intaractor := NewCustomerPropertyTypeGetAllInteractor(db.NewCustomerPropertyTypeRepository())
		response, err := intaractor.Handle()
		assert.NoError(t, err)

		////// 検証
		assert.Equal(t, preCreatedCustomerPropertyTypes, response.CustomerPropertyTypeDtos)
	})

	t.Run("データがなければ空スライスが返る", func(t *testing.T) {
		////// 準備
		// 既存データ削除
		conn, err := db.GetConnection()
		assert.NoError(t, err)
		defer conn.Db.Close()
		_, err = conn.Exec("DELETE FROM customer_types_customer_properties;")
		assert.NoError(t, err)
		_, err = conn.Exec("DELETE FROM customer_properties;")
		assert.NoError(t, err)

		////// 実行
		intaractor := NewCustomerPropertyTypeGetAllInteractor(db.NewCustomerPropertyTypeRepository())
		response, err := intaractor.Handle()
		assert.NoError(t, err)

		////// 検証
		assert.Len(t, response.CustomerPropertyTypeDtos, 0)
	})
}

func preCreateCustomerPropertyTypes() ([]customer_property_type.CustomerPropertyTypeDto, error) {
	retDtos := []customer_property_type.CustomerPropertyTypeDto{}
	intaractor := create.NewCustomerPropertyTypeCreateInteractor(db.NewCustomerPropertyTypeRepository())

	timestampStr := utils.CreateTimestampString()

	for i := 0; i < 10; i++ {
		num := strconv.Itoa(i)
		request := create.NewCustomerPropertyTypeCreateUseCaseRequest("プロパティ"+num+timestampStr, "string")
		response, err := intaractor.Handle(request)
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
