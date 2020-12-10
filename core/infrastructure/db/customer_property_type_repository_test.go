package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestCustomerPropertyTypeRepository_Create(t *testing.T) {
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
		savedId, err := r.Create(customerProperties, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証 todo: あとで再取得してデータが取れるか確認する
		assert.NotZero(t, savedId)
	})
}
