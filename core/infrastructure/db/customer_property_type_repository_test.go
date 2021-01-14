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

func TestCustomerPropertyTypeRepository_GetAll(t *testing.T) {
	rep := NewCustomerPropertyTypeRepository()

	conn, err := GetConnection()
	assert.NoError(t, err)
	defer conn.Db.Close()

	// トランザクション開始
	tran, err := conn.Begin()
	assert.NoError(t, err)

	////// 準備
	// 予め全データ削除
	_, err = tran.Exec("LOCK TABLE customer_types_customer_properties IN EXCLUSIVE MODE;")
	assert.NoError(t, err)
	_, err = tran.Exec("LOCK TABLE customer_properties IN EXCLUSIVE MODE;")
	assert.NoError(t, err)
	_, err = tran.Exec("DELETE FROM customers_customer_properties;")
	assert.NoError(t, err)
	_, err = tran.Exec("DELETE FROM customer_types_customer_properties;")
	assert.NoError(t, err)
	_, err = tran.Exec("DELETE FROM customer_properties;")
	assert.NoError(t, err)
	// データ新規登録
	timestampstr := utils.CreateTimestampString()
	customerProperties := []*customer.CustomerPropertyTypeEntity{
		customer.NewCustomerPropertyTypeEntity("性別"+timestampstr, customer.PROPERTY_TYPE_STRING),
		customer.NewCustomerPropertyTypeEntity("年齢"+timestampstr, customer.PROPERTY_TYPE_NUMERIC),
		customer.NewCustomerPropertyTypeEntity("住所"+timestampstr, customer.PROPERTY_TYPE_STRING),
	}
	// 事前データ保存実行
	savedIds, err := rep.Create(customerProperties, tran)
	assert.NoError(t, err)

	////// 実行
	loadedEntities, err := rep.GetAll(tran)
	assert.NoError(t, err)
	// トランザクションコミット
	err = tran.Commit()
	if err != nil {
		tran.Rollback()
		assert.Failf(t, "同時実行が影響してテストできなかった。err: %+v", err.Error())
	}

	////// 検証
	assert.Len(t, loadedEntities, len(customerProperties))
	for i, loadedEntity := range loadedEntities {
		assert.Equal(t, savedIds[i], loadedEntity.Id())
		assert.Equal(t, customerProperties[i].Name(), loadedEntity.Name())
		assert.Equal(t, customerProperties[i].PropertyType(), loadedEntity.PropertyType())
	}
}

func TestCustomerPropertyTypeRepository_GetByName(t *testing.T) {
	t.Run("Nameでデータを取得できる", func(t *testing.T) {
		////// 準備
		// 事前データ登録
		timestampstr := utils.CreateTimestampString()
		customerProperties := []*customer.CustomerPropertyTypeEntity{
			customer.NewCustomerPropertyTypeEntity("重複テスト"+timestampstr, customer.PROPERTY_TYPE_STRING),
		}
		conn, err := GetConnection()
		assert.NoError(t, err)
		tran, err := conn.Begin()
		assert.NoError(t, err)
		r := NewCustomerPropertyTypeRepository()
		savedIds, err := r.Create(customerProperties, tran)
		assert.NoError(t, err)
		err = tran.Commit()
		assert.NoError(t, err)

		////// 実行
		entity, err := r.GetByName("重複テスト"+timestampstr, conn)
		assert.NoError(t, err)

		////// 検証
		assert.Equal(t, savedIds[0], entity.Id())

	})
}
