package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestCustomerTypeRepository_Create_And_Get(t *testing.T) {
	t.Run("CustomerTypeエンティティを渡すとDBへ保存される", func(t *testing.T) {
		conn, err := GetConnection()
		assert.NoError(t, err)
		tran, err := conn.Begin()
		assert.NoError(t, err)

		////// 準備
		// カスタマープロパティタイプの登録
		timestamp := time.Now().UnixNano()
		timestampstr := strconv.Itoa(int(timestamp))
		customerProperties := []*customer.CustomerPropertyTypeEntity{
			customer.NewCustomerPropertyTypeEntity("性別"+timestampstr, customer.PROPERTY_TYPE_STRING),
			customer.NewCustomerPropertyTypeEntity("年齢"+timestampstr, customer.PROPERTY_TYPE_NUMERIC),
			customer.NewCustomerPropertyTypeEntity("住所"+timestampstr, customer.PROPERTY_TYPE_STRING),
		}
		propertyTypeRep := NewCustomerPropertyTypeRepository()
		savedPropertyIds, err := propertyTypeRep.Create(customerProperties, tran)
		assert.NoError(t, err)
		// カスタマータイプエンティティの作成
		customerType := customer.NewCustomerTypeEntity("顧客名"+timestampstr, savedPropertyIds)

		////// 実行
		r := NewCustomerTypeRepository()
		savedId, err := r.Create(customerType, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 検証 todo: あとで再取得してデータが取れるか確認する
		assert.NotZero(t, savedId)

		////// データの再取得
		reloadedEntity, err := r.GetById(savedId, conn)
		assert.NoError(t, err)

		////// 検証
		expected := customer.NewCustomerTypeEntityWithData(savedId, "顧客名"+timestampstr, savedPropertyIds)
		assert.True(t, reflect.DeepEqual(expected, reloadedEntity))
	})
}
