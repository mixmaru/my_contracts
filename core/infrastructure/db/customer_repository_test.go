package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/mixmaru/my_contracts/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomerRepository_Create(t *testing.T) {
	timestampStr := utils.CreateTimestampString()
	// カスタマータイプ作成
	kankouchoId, propertyIds, err := preCreateCustomerType(timestampStr)
	assert.NoError(t, err)

	t.Run("CustomerEntityを渡すと新規保存される", func(*testing.T) {
		////// 実行
		// カスタマーエンティティ作成
		newCustomer := customer.NewCustomerEntity(
			"厚生省"+timestampStr,
			kankouchoId,
			map[int]interface{}{
				propertyIds[0]: "03-1111-2222",
				propertyIds[1]: 200,
			},
		)
		rep := NewCustomerRepository()
		conn, err := GetConnection()
		assert.NoError(t, err)
		tran, err := conn.Begin()
		assert.NoError(t, err)
		savedId, err := rep.Create(newCustomer, tran)
		assert.NoError(t, err)
		assert.NotZero(t, savedId)

		////// 検証
		loadedEntity, err := rep.GetById(savedId, tran)
		assert.NotZero(t, loadedEntity.Id())
		assert.Equal(t, newCustomer.Name(), loadedEntity.Name())
		assert.Equal(t, newCustomer.CustomerTypeId(), loadedEntity.CustomerTypeId())
		assert.Equal(t, newCustomer.Properties(), loadedEntity.Properties())

		err = tran.Commit()
		assert.NoError(t, err)
	})

	t.Run("カスタマータイプで設定されていないプロパティを渡すとエラーになる", func(*testing.T) {
		////// 実行
		// カスタマーエンティティ作成
		newCustomer := customer.NewCustomerEntity(
			"厚生省"+timestampStr,
			kankouchoId,
			map[int]interface{}{
				propertyIds[0]: "03-1111-2222",
				propertyIds[1]: 200,
				192:            "他のカスタマータイプのプロパティ",
			},
		)
		rep := NewCustomerRepository()
		conn, err := GetConnection()
		assert.NoError(t, err)
		tran, err := conn.Begin()
		assert.NoError(t, err)
		savedId, err := rep.Create(newCustomer, tran)

		////// 検証
		assert.Error(t, err)
		assert.Zero(t, savedId)
		err = tran.Rollback()
		assert.NoError(t, err)
	})

	t.Run("カスタマータイプで設定されているプロパティが存在しないとそのプロパティは無視される", func(*testing.T) {

	})
}

func preCreateCustomerType(timestampStr string) (customerId int, propertyIds []int, err error) {
	// カスタマープロパティタイプ作成
	telId, err := preCreateCustomerProperty("電話番号"+timestampStr, customer.PROPERTY_TYPE_STRING)
	if err != nil {
		return 0, nil, err
	}
	countryId, err := preCreateCustomerProperty("何かしらの数値"+timestampStr, customer.PROPERTY_TYPE_NUMERIC)
	if err != nil {
		return 0, nil, err
	}
	propertyTypeIds := []int{
		telId,
		countryId,
	}

	// カスタマープロパティ作成
	newEntity := customer.NewCustomerTypeEntity("官公庁"+timestampStr, propertyTypeIds)
	rep := NewCustomerTypeRepository()
	conn, err := GetConnection()
	if err != nil {
		return 0, nil, err
	}
	defer conn.Db.Close()

	customerId, err = rep.Create(newEntity, conn)
	if err != nil {
		return 0, nil, err
	}
	propertyIds = []int{
		telId,
		countryId,
	}
	return customerId, propertyIds, nil
}

func preCreateCustomerProperty(name string, propertyType customer.PropertyType) (int, error) {
	newEntity := customer.NewCustomerPropertyTypeEntity(name, propertyType)
	rep := NewCustomerPropertyTypeRepository()

	conn, err := GetConnection()
	if err != nil {
		return 0, err
	}
	savedIds, err := rep.Create([]*customer.CustomerPropertyTypeEntity{newEntity}, conn)
	if err != nil {
		return 0, err
	}
	return savedIds[0], nil
}
