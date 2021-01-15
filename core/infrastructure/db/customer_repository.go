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
	err = crateCustomerCustomerProperties(savedCustomerId, customerEntity.CustomerTypeId(), customerEntity.Properties(), executor)
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

func crateCustomerCustomerProperties(customerId, customerTypeId int, properties map[int]interface{}, executor gorp.SqlExecutor) error {
	// mapper作成
	mappers := make([]interface{}, 0, len(properties))
	for key, val := range properties {
		value, err := toText(val)
		if err != nil {
			return err
		}
		mapper := customerCustomerPropertyMapper{
			CustomerId:         customerId,
			CustomerTypeId:     customerTypeId,
			CustomerPropertyId: key,
			Value:              value,
		}
		mappers = append(mappers, &mapper)
	}
	// 保存実行
	err := executor.Insert(mappers...)
	if err != nil {
		return errors.Wrapf(err, "customerPropertiesの保存に失敗しました。mappers: %+v", mappers...)
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
		return "", errors.Errorf("string型へ変換できなかった。value: %+v, value.(type): %T", value, value)
	}
}

func (c *CustomerRepository) GetById(id int, executor gorp.SqlExecutor) (entity *customer.CustomerEntity, err error) {
	query := `
select
       c.id,
       c.customer_type_id,
       c.name,
       ccp.customer_property_id,
       cp.type,
       ccp.value
from customers c
inner join customers_customer_properties ccp on c.id = ccp.customer_id
inner join customer_properties cp on ccp.customer_property_id = cp.id
inner join customer_types_customer_properties ctcp on cp.id = ctcp.customer_property_id
where c.id = $1
order by ctcp."order"
`
	var mappers []*customerMapperForGet

	_, err = executor.Select(&mappers, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "データ取得失敗. query: %+v, id: %v", query, id)
	}

	retEntity, err := generateEntity(mappers)
	if err != nil {
		return nil, err
	}
	return retEntity, nil
}

type customerMapperForGet struct {
	Id                 int                   `db:"id"`
	CustomerTypeId     int                   `db:"customer_type_id"`
	Name               string                `db:"name"`
	CustomerPropertyId int                   `db:"customer_property_id"`
	PropertyType       customer.PropertyType `db:"type"`
	Value              string                `db:value`
}

func generateEntity(mappers []*customerMapperForGet) (*customer.CustomerEntity, error) {
	properties := map[int]interface{}{}
	for _, mapper := range mappers {
		var err error
		properties[mapper.CustomerPropertyId], err = valueFromText(mapper.PropertyType, mapper.Value)
		if err != nil {
			return nil, err
		}
	}
	entity := customer.NewCustomerEntityWithData(
		mappers[0].Id,
		mappers[0].Name,
		mappers[0].CustomerTypeId,
		properties,
	)
	return entity, nil
}

func valueFromText(propertyType customer.PropertyType, value string) (interface{}, error) {
	switch propertyType {
	case customer.PROPERTY_TYPE_STRING:
		return value, nil
	case customer.PROPERTY_TYPE_NUMERIC:
		retValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, errors.Wrapf(err, "intへの変換に失敗した。value: %+v", value)
		}
		return retValue, nil
	default:
		return nil, errors.Errorf("想定外のエラー")
	}
}

type customerMapper struct {
	Id             int    `db:"id"`
	Name           string `db:"name"`
	CustomerTypeId int    `db:"customer_type_id"`
	CreatedAtUpdatedAtMapper
}

type customerCustomerPropertyMapper struct {
	CustomerId         int    `db:"customer_id"`
	CustomerTypeId     int    `db:"customer_type_id"`
	CustomerPropertyId int    `db:"customer_property_id"`
	Value              string `db:"value"`
	CreatedAtUpdatedAtMapper
}
