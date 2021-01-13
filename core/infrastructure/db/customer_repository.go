package db

import (
	"github.com/mixmaru/my_contracts/core/domain/models/customer"
	"github.com/pkg/errors"
	"gopkg.in/gorp.v2"
	"strconv"
)

type CustomerRepository struct {
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{}
}

// カスタマー新規作成
func (c *CustomerRepository) Create(customerEntity *customer.CustomerEntity, executor gorp.SqlExecutor) (savedId int, err error) {
	// カスタマー登録
	savedCustomerId, err := createCustomer(customerEntity, executor)
	if err != nil {
		return 0, err
	}

	// カスタマープロパティ登録
	err = crateCustomerCustomerProperties(savedCustomerId, customerEntity.Properties(), executor)
	if err != nil {
		return 0, err
	}

	return savedCustomerId, nil
}

func createCustomer(customerEntity *customer.CustomerEntity, executor gorp.SqlExecutor) (int, error) {
	// mapper作成
	newCustomer := customerMapper{
		Name:           customerEntity.Name(),
		CustomerTypeId: customerEntity.CustomerTypeId(),
	}
	// 保存実行
	err := executor.Insert(&newCustomer)
	if err != nil {
		return 0, errors.Wrapf(err, "Customerテーブルへのデータ保存エラー。customerEntity: %+v", customerEntity)
	}
	return newCustomer.Id, nil
}

func crateCustomerCustomerProperties(customerId int, properties map[int]interface{}, executor gorp.SqlExecutor) error {
	// mapper作成
	mappers := make([]interface{}, 0, len(properties))
	for key, val := range properties {
		value, err := toText(val)
		if err != nil {
			return err
		}
		mapper := customerCustomerPropertyMapper{
			CustomerId:         customerId,
			CustomerPropertyId: key,
			Value:              value,
		}
		mappers = append(mappers, &mapper)
	}
	// 保存実行
	err := executor.Insert(mappers...)
	if err != nil {
		return errors.Wrapf(err, "customerPropertiesの保存に失敗しました。mappers: %+v", mappers)
	}
	return nil
}

// value(int or string方のinterface{}型)をstringに変換する
func toText(value interface{}) (string, error) {
	switch value.(type) {
	case string:
		return value.(string), nil
	case int:
		return strconv.Itoa(value.(int)), nil
	default:
		return "", errors.Errorf("string型へ変換できなかった。value.(type): %T", value)
	}
}

func (c *CustomerRepository) GetById(id int, executor gorp.SqlExecutor) (entity *customer.CustomerEntity, err error) {
	return nil, nil
}

type customerMapper struct {
	Id             int    `db:"id"`
	Name           string `db:"name"`
	CustomerTypeId int    `db:"customer_type_id"`
	CreatedAtUpdatedAtMapper
}

type customerCustomerPropertyMapper struct {
	CustomerId         int    `db:"customer_id"`
	CustomerPropertyId int    `db:"customer_property_id"`
	Value              string `db:"value"`
	CreatedAtUpdatedAtMapper
}
