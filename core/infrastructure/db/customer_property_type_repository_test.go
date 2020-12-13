package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCustomerPropertyTypeRepository_Create_And_GetByIds(t *testing.T) {
	t.Run("CustomerPropertyTypeエンティティを渡すとDBへ保存される", func(t *testing.T) {
		////// 準備
		timestampstr := utils.CreateTimestampString()
		customerProperties := []*customer.CustomerPropertyTypeEntity{
			customer.NewCustomerPropertyTypeEntity("性別"+timestampstr, customer.PROPERTY_TYPE_STRING),
			customer.NewCustomerPropertyTypeEntity("年齢"+timestampstr, customer.PROPERTY_TYPE_NUMERIC),
			customer.NewCustomerPropertyTypeEntity("住所"+timestampstr, customer.PROPERTY_TYPE_STRING),
		}

		////// 実行
		conn, err := GetConnection()
		assert.NoError(t, err)
		tran, err := conn.Begin()
		assert.NoError(t, err)
		r := NewCustomerPropertyTypeRepository()
		savedIds, err := r.Create(customerProperties, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)
		assert.NotZero(t, savedIds)

		////// 保存したデータを取得してみる
		actual, err := r.GetByIds(savedIds, conn)
		assert.NoError(t, err)
		////// 検証
		expected := []*customer.CustomerPropertyTypeEntity{
			customer.NewCustomerPropertyTypeEntityWithData(savedIds[0], "性別"+timestampstr, customer.PROPERTY_TYPE_STRING),
			customer.NewCustomerPropertyTypeEntityWithData(savedIds[1], "年齢"+timestampstr, customer.PROPERTY_TYPE_NUMERIC),
			customer.NewCustomerPropertyTypeEntityWithData(savedIds[2], "住所"+timestampstr, customer.PROPERTY_TYPE_STRING),
		}
		assert.True(t, reflect.DeepEqual(actual, expected))
	})
}
