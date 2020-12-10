package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestCustomerPropertyTypeRepository_Create_And_GetByIds(t *testing.T) {
	t.Run("CustomerPropertyTypeエンティティを渡すとDBへ保存される", func(t *testing.T) {
		////// 準備
		timestamp := time.Now().UnixNano()
		timestampstr := strconv.Itoa(int(timestamp))
		customerProperties := []*customer.CustomerPropertyTypeEntity{
			customer.NewCustomerParamTypeEntity("性別"+timestampstr, customer.PROPERTY_TYPE_STRING),
			customer.NewCustomerParamTypeEntity("年齢"+timestampstr, customer.PROPERTY_TYPE_NUMERIC),
			customer.NewCustomerParamTypeEntity("住所"+timestampstr, customer.PROPERTY_TYPE_STRING),
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
			customer.NewCustomerParamTypeEntityWithData(savedIds[0], "性別"+timestampstr, customer.PROPERTY_TYPE_STRING),
			customer.NewCustomerParamTypeEntityWithData(savedIds[1], "年齢"+timestampstr, customer.PROPERTY_TYPE_NUMERIC),
			customer.NewCustomerParamTypeEntityWithData(savedIds[2], "住所"+timestampstr, customer.PROPERTY_TYPE_STRING),
		}
		assert.True(t, reflect.DeepEqual(actual, expected))
	})
}
